package hive_dirver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"gorm.io/gorm/migrator"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	DriverName = "hive"
)

type Config struct {
	DriverName    string
	ServerVersion string
	DSN           string
	DSNConfig     *DSNConfig
	Conn          gorm.ConnPool

	DefaultStringSize uint
}

type Dialector struct {
	*Config

	logger logger.Interface
}

var (
	// CreateClauses create clauses
	CreateClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	// QueryClauses query clauses
	QueryClauses = []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR"}
	// UpdateClauses update clauses
	UpdateClauses = []string{"UPDATE", "SET", "WHERE"}
	// DeleteClauses delete clauses
	DeleteClauses = []string{"DELETE", "FROM", "WHERE"}
)

func Open(dsn string) gorm.Dialector {
	dsnConf, _ := ParseDSN(dsn)
	dsn = dsnConf.Complete().FormatDSN()
	return &Dialector{Config: &Config{DSN: dsn, DSNConfig: dsnConf}}
}

func New(config Config) gorm.Dialector {
	switch {
	case config.DSN == "" && config.DSNConfig != nil:
		config.DSN = config.DSNConfig.Complete().FormatDSN()
	case config.DSN != "" && config.DSNConfig == nil:
		config.DSNConfig, _ = ParseDSN(config.DSN)
		config.DSNConfig.Complete()
	}
	return &Dialector{Config: &config}
}

func (dialector *Dialector) Name() string {
	return DriverName
}

func (dialector *Dialector) Initialize(db *gorm.DB) (err error) {
	dialector.logger = db.Logger
	if dialector.DriverName == "" {
		dialector.DriverName = DriverName
	}

	if dialector.Conn != nil {
		db.ConnPool = dialector.Conn
	} else {
		dsn := dialector.DSNConfig.FormatDSN()
		db.ConnPool, err = sql.Open(dialector.DriverName, dsn)
		if err != nil {
			return err
		}
	}

	callbackConfig := &callbacks.Config{
		CreateClauses: CreateClauses,
		QueryClauses:  QueryClauses,
		UpdateClauses: UpdateClauses,
		DeleteClauses: DeleteClauses,
	}

	dialector.RegisterCallbacks(db, callbackConfig)

	if db.Config != nil {
		db.Config.PrepareStmt = false
		db.Config.DisableNestedTransaction = false
		db.Config.SkipDefaultTransaction = true
	}

	return nil
}

func (dialector *Dialector) RegisterCallbacks(db *gorm.DB, config *callbacks.Config) {

	if len(config.CreateClauses) == 0 {
		config.CreateClauses = CreateClauses
	}
	if len(config.QueryClauses) == 0 {
		config.QueryClauses = QueryClauses
	}
	if len(config.DeleteClauses) == 0 {
		config.DeleteClauses = DeleteClauses
	}
	if len(config.UpdateClauses) == 0 {
		config.UpdateClauses = UpdateClauses
	}

	createCallback := db.Callback().Create()
	dialector.handleError(createCallback.Register("gorm:before_create", callbacks.BeforeCreate))
	dialector.handleError(createCallback.Register("gorm:save_before_associations", callbacks.SaveBeforeAssociations(true)))
	dialector.handleError(createCallback.Register("gorm:create", callbacks.Create(config)))
	dialector.handleError(createCallback.Register("gorm:save_after_associations", callbacks.SaveAfterAssociations(true)))
	dialector.handleError(createCallback.Register("gorm:after_create", callbacks.AfterCreate))
	createCallback.Clauses = config.CreateClauses

	queryCallback := db.Callback().Query()
	dialector.handleError(queryCallback.Register("gorm:query", callbacks.Query))
	dialector.handleError(queryCallback.Register("gorm:preload", callbacks.Preload))
	dialector.handleError(queryCallback.Register("gorm:after_query", callbacks.AfterQuery))
	queryCallback.Clauses = config.QueryClauses

	deleteCallback := db.Callback().Delete()
	dialector.handleError(deleteCallback.Register("gorm:before_delete", callbacks.BeforeDelete))
	dialector.handleError(deleteCallback.Register("gorm:delete_before_associations", callbacks.DeleteBeforeAssociations))
	dialector.handleError(deleteCallback.Register("gorm:delete", callbacks.Delete(config)))
	dialector.handleError(deleteCallback.Register("gorm:after_delete", callbacks.AfterDelete))
	deleteCallback.Clauses = config.DeleteClauses

	updateCallback := db.Callback().Update()
	dialector.handleError(updateCallback.Register("gorm:setup_reflect_value", callbacks.SetupUpdateReflectValue))
	dialector.handleError(updateCallback.Register("gorm:before_update", callbacks.BeforeUpdate))
	dialector.handleError(updateCallback.Register("gorm:save_before_associations", callbacks.SaveBeforeAssociations(false)))
	dialector.handleError(updateCallback.Register("gorm:update", callbacks.Update(config)))
	dialector.handleError(updateCallback.Register("gorm:save_after_associations", callbacks.SaveAfterAssociations(false)))
	dialector.handleError(updateCallback.Register("gorm:after_update", callbacks.AfterUpdate))
	updateCallback.Clauses = config.UpdateClauses

	rowCallback := db.Callback().Row()
	dialector.handleError(rowCallback.Register("gorm:row", callbacks.RowQuery))
	rowCallback.Clauses = config.QueryClauses

	rawCallback := db.Callback().Raw()
	dialector.handleError(rawCallback.Register("gorm:raw", callbacks.RawExec))
	rawCallback.Clauses = config.QueryClauses
}

func (dialector *Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{
		Migrator: migrator.Migrator{
			Config: migrator.Config{
				DB:        db,
				Dialector: dialector,
			},
		},
		Dialector: dialector,
	}
}

func (dialector *Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "boolean"
	case schema.Int, schema.Uint:
		return dialector.getSchemaIntAndUnitType(field)
	case schema.Float:
		return dialector.getSchemaFloatType(field)
	case schema.String:
		return dialector.getSchemaStringType(field)
	case schema.Time:
		return dialector.getSchemaTimeType(field)
	case schema.Bytes:
		return dialector.getSchemaBytesType(field)
	default:
		return dialector.getSchemaCustomType(field)
	}
}

func (dialector *Dialector) getSchemaIntAndUnitType(field *schema.Field) string {
	switch {
	case field.Size <= 8:
		return "tinyint"
	case field.Size <= 16:
		return "smallint"
	case field.Size <= 32:
		return "int"
	default:
		return "bigint"
	}
}

func (dialector *Dialector) getSchemaFloatType(field *schema.Field) string {
	if field.Precision > 0 {
		return fmt.Sprintf("decimal(%d, %d)", field.Precision, field.Scale)
	}

	if field.Size <= 32 {
		return "float"
	}

	return "double"
}

func (dialector *Dialector) getSchemaStringType(field *schema.Field) string {
	return "String" // TODO: varchar?
}

func (dialector *Dialector) getSchemaTimeType(field *schema.Field) string {
	return "Timestamp" // TODO: DATE?
}

func (dialector *Dialector) getSchemaBytesType(field *schema.Field) string {
	return "BINARY"
}

func (dialector *Dialector) getSchemaCustomType(field *schema.Field) string {
	sqlType := string(field.DataType)
	return sqlType
}

func (dialector *Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: "DEFAULT"}
}

func (dialector *Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	_ = writer.WriteByte('?')
}

func (dialector *Dialector) QuoteTo(writer clause.Writer, str string) {
	_ = writer.WriteByte('`')
	if strings.Contains(str, ".") {
		for idx, str := range strings.Split(str, ".") {
			if idx > 0 {
				_, _ = writer.WriteString(".`")
			}
			_, _ = writer.WriteString(str)
			_ = writer.WriteByte('`')
		}
	} else {
		_, _ = writer.WriteString(str)
		_ = writer.WriteByte('`')
	}
}

func (dialector *Dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `'`, vars...)
}

func (dialector *Dialector) getLogger() logger.Interface {
	if dialector.logger == nil {
		return logger.Default
	}
	return dialector.logger
}

func (dialector *Dialector) handleError(err error, ignoreErrors ...error) {
	if err != nil {
		for _, except := range ignoreErrors {
			if errors.Is(err, except) {
				return
			}
		}
		_, file, line, _ := runtime.Caller(1)
		dialector.getLogger().Warn(context.Background(), "%s:%d %v", file, line, err)
	}
}
