//节点配置层
package es_settings

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic"
)

type Setting struct {
	Persistent struct {
		Action struct {
			AutoCreateIndex string `json:"auto_create_index"`
		} `json:"action"`
		Xpack struct {
			Monitoring struct {
				Collection struct {
					Enabled string `json:"enabled"`
				} `json:"collection"`
			} `json:"monitoring"`
		} `json:"xpack"`
	} `json:"persistent"`
	Transient struct {
	} `json:"transient"`
	Defaults struct {
		Cluster struct {
			Routing struct {
				UseAdaptiveReplicaSelection string `json:"use_adaptive_replica_selection"`
				Rebalance                   struct {
					Enable string `json:"enable"`
				} `json:"rebalance"`
				Allocation struct {
					NodeConcurrentIncomingRecoveries string `json:"node_concurrent_incoming_recoveries"`
					NodeInitialPrimariesRecoveries   string `json:"node_initial_primaries_recoveries"`
					SameShard                        struct {
						Host string `json:"host"`
					} `json:"same_shard"`
					TotalShardsPerNode string `json:"total_shards_per_node"`
					Type               string `json:"type"`
					Disk               struct {
						ThresholdEnabled string `json:"threshold_enabled"`
						Watermark        struct {
							Low        string `json:"low"`
							FloodStage string `json:"flood_stage"`
							High       string `json:"high"`
						} `json:"watermark"`
						IncludeRelocations string `json:"include_relocations"`
						RerouteInterval    string `json:"reroute_interval"`
					} `json:"disk"`
					Awareness struct {
						Attributes []interface{} `json:"attributes"`
					} `json:"awareness"`
					Balance struct {
						Index     string `json:"index"`
						Threshold string `json:"threshold"`
						Shard     string `json:"shard"`
					} `json:"balance"`
					Enable                           string `json:"enable"`
					NodeConcurrentOutgoingRecoveries string `json:"node_concurrent_outgoing_recoveries"`
					AllowRebalance                   string `json:"allow_rebalance"`
					ClusterConcurrentRebalance       string `json:"cluster_concurrent_rebalance"`
					NodeConcurrentRecoveries         string `json:"node_concurrent_recoveries"`
				} `json:"allocation"`
			} `json:"routing"`
			Indices struct {
				Tombstones struct {
					Size string `json:"size"`
				} `json:"tombstones"`
				Close struct {
					Enable string `json:"enable"`
				} `json:"close"`
			} `json:"indices"`
			Nodes struct {
				ReconnectInterval string `json:"reconnect_interval"`
			} `json:"nodes"`
			PersistentTasks struct {
				Allocation struct {
					Enable string `json:"enable"`
				} `json:"allocation"`
			} `json:"persistent_tasks"`
			Blocks struct {
				ReadOnlyAllowDelete string `json:"read_only_allow_delete"`
				ReadOnly            string `json:"read_only"`
			} `json:"blocks"`
			Service struct {
				SlowTaskLoggingThreshold string `json:"slow_task_logging_threshold"`
			} `json:"service"`
			Name             string `json:"name"`
			MaxShardsPerNode string `json:"max_shards_per_node"`
			Remote           struct {
				Node struct {
					Attr string `json:"attr"`
				} `json:"node"`
				InitialConnectTimeout string `json:"initial_connect_timeout"`
				Connect               string `json:"connect"`
				ConnectionsPerCluster string `json:"connections_per_cluster"`
			} `json:"remote"`
			Info struct {
				Update struct {
					Interval string `json:"interval"`
					Timeout  string `json:"timeout"`
				} `json:"update"`
			} `json:"info"`
		} `json:"cluster"`
		No struct {
			Model struct {
				State struct {
					Persist string `json:"persist"`
				} `json:"state"`
			} `json:"model"`
		} `json:"no"`
		Logger struct {
			Level string `json:"level"`
		} `json:"logger"`
		Bootstrap struct {
			MemoryLock       string `json:"memory_lock"`
			SystemCallFilter string `json:"system_call_filter"`
			Ctrlhandler      string `json:"ctrlhandler"`
		} `json:"bootstrap"`
		Processors string `json:"processors"`
		Ingest     struct {
			Grok struct {
				Watchdog struct {
					MaxExecutionTime string `json:"max_execution_time"`
					Interval         string `json:"interval"`
				} `json:"watchdog"`
			} `json:"grok"`
		} `json:"ingest"`
		Network struct {
			Host []string `json:"host"`
			TCP  struct {
				ReuseAddress      string `json:"reuse_address"`
				KeepAlive         string `json:"keep_alive"`
				ConnectTimeout    string `json:"connect_timeout"`
				ReceiveBufferSize string `json:"receive_buffer_size"`
				NoDelay           string `json:"no_delay"`
				SendBufferSize    string `json:"send_buffer_size"`
			} `json:"tcp"`
			BindHost []string `json:"bind_host"`
			Server   string   `json:"server"`
			Breaker  struct {
				InflightRequests struct {
					Limit    string `json:"limit"`
					Overhead string `json:"overhead"`
				} `json:"inflight_requests"`
			} `json:"breaker"`
			PublishHost []string `json:"publish_host"`
		} `json:"network"`
		Pidfile string `json:"pidfile"`
		Path    struct {
			Data       []interface{} `json:"data"`
			Logs       string        `json:"logs"`
			SharedData string        `json:"shared_data"`
			Home       string        `json:"home"`
			Repo       []string      `json:"repo"`
		} `json:"path"`
		Search struct {
			DefaultSearchTimeout string `json:"default_search_timeout"`
			Highlight            struct {
				TermVectorMultiValue string `json:"term_vector_multi_value"`
			} `json:"highlight"`
			DefaultAllowPartialResults string `json:"default_allow_partial_results"`
			MaxBuckets                 string `json:"max_buckets"`
			LowLevelCancellation       string `json:"low_level_cancellation"`
			KeepAliveInterval          string `json:"keep_alive_interval"`
			Remote                     struct {
				Node struct {
					Attr string `json:"attr"`
				} `json:"node"`
				InitialConnectTimeout string `json:"initial_connect_timeout"`
				Connect               string `json:"connect"`
				ConnectionsPerCluster string `json:"connections_per_cluster"`
			} `json:"remote"`
			DefaultKeepAlive string `json:"default_keep_alive"`
			MaxKeepAlive     string `json:"max_keep_alive"`
		} `json:"search"`
		Security struct {
			Manager struct {
				FilterBadDefaults string `json:"filter_bad_defaults"`
			} `json:"manager"`
		} `json:"security"`
		Repositories struct {
			Fs struct {
				Compress  string `json:"compress"`
				ChunkSize string `json:"chunk_size"`
				Location  string `json:"location"`
			} `json:"fs"`
			URL struct {
				SupportedProtocols []string      `json:"supported_protocols"`
				AllowedUrls        []interface{} `json:"allowed_urls"`
				URL                string        `json:"url"`
			} `json:"url"`
		} `json:"repositories"`
		Action struct {
			Search struct {
				ShardCount struct {
					Limit string `json:"limit"`
				} `json:"shard_count"`
			} `json:"search"`
			DestructiveRequiresName string `json:"destructive_requires_name"`
			Master                  struct {
				ForceLocal string `json:"force_local"`
			} `json:"master"`
		} `json:"action"`
		Client struct {
			Type      string `json:"type"`
			Transport struct {
				IgnoreClusterName    string `json:"ignore_cluster_name"`
				NodesSamplerInterval string `json:"nodes_sampler_interval"`
				Sniff                string `json:"sniff"`
				PingTimeout          string `json:"ping_timeout"`
			} `json:"transport"`
		} `json:"client"`
		Xpack struct {
			Watcher struct {
				Execution struct {
					Scroll struct {
						Size    string `json:"size"`
						Timeout string `json:"timeout"`
					} `json:"scroll"`
					DefaultThrottlePeriod string `json:"default_throttle_period"`
				} `json:"execution"`
				Internal struct {
					Ops struct {
						Bulk struct {
							DefaultTimeout string `json:"default_timeout"`
						} `json:"bulk"`
						Index struct {
							DefaultTimeout string `json:"default_timeout"`
						} `json:"index"`
						Search struct {
							DefaultTimeout string `json:"default_timeout"`
						} `json:"search"`
					} `json:"ops"`
				} `json:"internal"`
				ThreadPool struct {
					QueueSize string `json:"queue_size"`
					Size      string `json:"size"`
				} `json:"thread_pool"`
				Index struct {
					Rest struct {
						DirectAccess string `json:"direct_access"`
					} `json:"rest"`
				} `json:"index"`
				History struct {
					CleanerService struct {
						Enabled string `json:"enabled"`
					} `json:"cleaner_service"`
				} `json:"history"`
				Trigger struct {
					Schedule struct {
						Ticker struct {
							TickInterval string `json:"tick_interval"`
						} `json:"ticker"`
					} `json:"schedule"`
				} `json:"trigger"`
				Enabled string `json:"enabled"`
				Input   struct {
					Search struct {
						DefaultTimeout string `json:"default_timeout"`
					} `json:"search"`
				} `json:"input"`
				EncryptSensitiveData string `json:"encrypt_sensitive_data"`
				Transform            struct {
					Search struct {
						DefaultTimeout string `json:"default_timeout"`
					} `json:"search"`
				} `json:"transform"`
				Stop struct {
					Timeout string `json:"timeout"`
				} `json:"stop"`
				Watch struct {
					Scroll struct {
						Size string `json:"size"`
					} `json:"scroll"`
				} `json:"watch"`
				RequireManualStart string `json:"require_manual_start"`
				Bulk               struct {
					ConcurrentRequests string `json:"concurrent_requests"`
					FlushInterval      string `json:"flush_interval"`
					Size               string `json:"size"`
					Actions            string `json:"actions"`
				} `json:"bulk"`
				Actions struct {
					Bulk struct {
						DefaultTimeout string `json:"default_timeout"`
					} `json:"bulk"`
					Index struct {
						DefaultTimeout string `json:"default_timeout"`
					} `json:"index"`
				} `json:"actions"`
			} `json:"watcher"`
			License struct {
				SelfGenerated struct {
					Type string `json:"type"`
				} `json:"self_generated"`
			} `json:"license"`
			Logstash struct {
				Enabled string `json:"enabled"`
			} `json:"logstash"`
			Notification struct {
				Hipchat struct {
					Host           string `json:"host"`
					Port           string `json:"port"`
					DefaultAccount string `json:"default_account"`
				} `json:"hipchat"`
				Pagerduty struct {
					DefaultAccount string `json:"default_account"`
				} `json:"pagerduty"`
				Email struct {
					DefaultAccount string `json:"default_account"`
					HTML           struct {
						Sanitization struct {
							Allow    []string      `json:"allow"`
							Disallow []interface{} `json:"disallow"`
							Enabled  string        `json:"enabled"`
						} `json:"sanitization"`
					} `json:"html"`
				} `json:"email"`
				Reporting struct {
					Retries  string `json:"retries"`
					Interval string `json:"interval"`
				} `json:"reporting"`
				Jira struct {
					DefaultAccount string `json:"default_account"`
				} `json:"jira"`
				Slack struct {
					DefaultAccount string `json:"default_account"`
				} `json:"slack"`
			} `json:"notification"`
			Security struct {
				DlsFls struct {
					Enabled string `json:"enabled"`
				} `json:"dls_fls"`
				Transport struct {
					Filter struct {
						Allow   []interface{} `json:"allow"`
						Deny    []interface{} `json:"deny"`
						Enabled string        `json:"enabled"`
					} `json:"filter"`
					Ssl struct {
						Enabled string `json:"enabled"`
					} `json:"ssl"`
				} `json:"transport"`
				Enabled string `json:"enabled"`
				Filter  struct {
					AlwaysAllowBoundAddress string `json:"always_allow_bound_address"`
				} `json:"filter"`
				Encryption struct {
					Algorithm string `json:"algorithm"`
				} `json:"encryption"`
				Audit struct {
					Outputs []string `json:"outputs"`
					Index   struct {
						BulkSize      string `json:"bulk_size"`
						Rollover      string `json:"rollover"`
						FlushInterval string `json:"flush_interval"`
						Events        struct {
							EmitRequestBody string        `json:"emit_request_body"`
							Include         []string      `json:"include"`
							Exclude         []interface{} `json:"exclude"`
						} `json:"events"`
						QueueMaxSize string `json:"queue_max_size"`
					} `json:"index"`
					Enabled string `json:"enabled"`
					Logfile struct {
						EmitNodeID       string `json:"emit_node_id"`
						EmitNodeHostName string `json:"emit_node_host_name"`
						EmitNodeName     string `json:"emit_node_name"`
						Events           struct {
							EmitRequestBody string        `json:"emit_request_body"`
							Include         []string      `json:"include"`
							Exclude         []interface{} `json:"exclude"`
						} `json:"events"`
						EmitNodeHostAddress string `json:"emit_node_host_address"`
					} `json:"logfile"`
				} `json:"audit"`
				Authc struct {
					Anonymous struct {
						AuthzException string        `json:"authz_exception"`
						Roles          []interface{} `json:"roles"`
						Username       string        `json:"username"`
					} `json:"anonymous"`
					PasswordHashing struct {
						Algorithm string `json:"algorithm"`
					} `json:"password_hashing"`
					RunAs struct {
						Enabled string `json:"enabled"`
					} `json:"run_as"`
					ReservedRealm struct {
						Enabled string `json:"enabled"`
					} `json:"reserved_realm"`
					Token struct {
						Delete struct {
							Interval string `json:"interval"`
							Timeout  string `json:"timeout"`
						} `json:"delete"`
						Enabled    string `json:"enabled"`
						ThreadPool struct {
							QueueSize string `json:"queue_size"`
							Size      string `json:"size"`
						} `json:"thread_pool"`
						Timeout string `json:"timeout"`
					} `json:"token"`
				} `json:"authc"`
				FipsMode struct {
					Enabled string `json:"enabled"`
				} `json:"fips_mode"`
				EncryptionKey struct {
					Length    string `json:"length"`
					Algorithm string `json:"algorithm"`
				} `json:"encryption_key"`
				HTTP struct {
					Filter struct {
						Allow   []interface{} `json:"allow"`
						Deny    []interface{} `json:"deny"`
						Enabled string        `json:"enabled"`
					} `json:"filter"`
					Ssl struct {
						Enabled string `json:"enabled"`
					} `json:"ssl"`
				} `json:"http"`
				Automata struct {
					MaxDeterminizedStates string `json:"max_determinized_states"`
					Cache                 struct {
						Size    string `json:"size"`
						TTL     string `json:"ttl"`
						Enabled string `json:"enabled"`
					} `json:"cache"`
				} `json:"automata"`
				User  interface{} `json:"user"`
				Authz struct {
					Store struct {
						Roles struct {
							Index struct {
								Cache struct {
									TTL     string `json:"ttl"`
									MaxSize string `json:"max_size"`
								} `json:"cache"`
							} `json:"index"`
							Cache struct {
								MaxSize string `json:"max_size"`
							} `json:"cache"`
							NegativeLookupCache struct {
								MaxSize string `json:"max_size"`
							} `json:"negative_lookup_cache"`
							FieldPermissions struct {
								Cache struct {
									MaxSizeInBytes string `json:"max_size_in_bytes"`
								} `json:"cache"`
							} `json:"field_permissions"`
						} `json:"roles"`
					} `json:"store"`
				} `json:"authz"`
			} `json:"security"`
			Ccr struct {
				Enabled    string `json:"enabled"`
				AutoFollow struct {
					PollInterval string `json:"poll_interval"`
				} `json:"auto_follow"`
				CcrThreadPool struct {
					QueueSize string `json:"queue_size"`
					Size      string `json:"size"`
				} `json:"ccr_thread_pool"`
			} `json:"ccr"`
			HTTP struct {
				DefaultConnectionTimeout string `json:"default_connection_timeout"`
				Proxy                    struct {
					Host   string `json:"host"`
					Scheme string `json:"scheme"`
					Port   string `json:"port"`
				} `json:"proxy"`
				DefaultReadTimeout string `json:"default_read_timeout"`
				MaxResponseSize    string `json:"max_response_size"`
			} `json:"http"`
			Monitoring struct {
				Collection struct {
					Cluster struct {
						Stats struct {
							Timeout string `json:"timeout"`
						} `json:"stats"`
					} `json:"cluster"`
					Node struct {
						Stats struct {
							Timeout string `json:"timeout"`
						} `json:"stats"`
					} `json:"node"`
					Indices []interface{} `json:"indices"`
					Ccr     struct {
						Stats struct {
							Timeout string `json:"timeout"`
						} `json:"stats"`
					} `json:"ccr"`
					Index struct {
						Stats struct {
							Timeout string `json:"timeout"`
						} `json:"stats"`
						Recovery struct {
							ActiveOnly string `json:"active_only"`
							Timeout    string `json:"timeout"`
						} `json:"recovery"`
					} `json:"index"`
					Interval string `json:"interval"`
					Ml       struct {
						Job struct {
							Stats struct {
								Timeout string `json:"timeout"`
							} `json:"stats"`
						} `json:"job"`
					} `json:"ml"`
				} `json:"collection"`
				History struct {
					Duration string `json:"duration"`
				} `json:"history"`
				Elasticsearch struct {
					Collection struct {
						Enabled string `json:"enabled"`
					} `json:"collection"`
				} `json:"elasticsearch"`
				Enabled string `json:"enabled"`
			} `json:"monitoring"`
			Graph struct {
				Enabled string `json:"enabled"`
			} `json:"graph"`
			Ml struct {
				UtilityThreadPool struct {
					QueueSize string `json:"queue_size"`
					Size      string `json:"size"`
				} `json:"utility_thread_pool"`
				MaxAnomalyRecords       string `json:"max_anomaly_records"`
				MaxMachineMemoryPercent string `json:"max_machine_memory_percent"`
				MaxOpenJobs             string `json:"max_open_jobs"`
				MinDiskSpaceOffHeap     string `json:"min_disk_space_off_heap"`
				AutodetectProcess       string `json:"autodetect_process"`
				DatafeedThreadPool      struct {
					QueueSize string `json:"queue_size"`
					Size      string `json:"size"`
				} `json:"datafeed_thread_pool"`
				NodeConcurrentJobAllocations string `json:"node_concurrent_job_allocations"`
				MaxModelMemoryLimit          string `json:"max_model_memory_limit"`
				Enabled                      string `json:"enabled"`
				MaxLazyMlNodes               string `json:"max_lazy_ml_nodes"`
				AutodetectThreadPool         struct {
					QueueSize string `json:"queue_size"`
					Size      string `json:"size"`
				} `json:"autodetect_thread_pool"`
			} `json:"ml"`
			Rollup struct {
				Enabled        string `json:"enabled"`
				TaskThreadPool struct {
					QueueSize string `json:"queue_size"`
					Size      string `json:"size"`
				} `json:"task_thread_pool"`
			} `json:"rollup"`
			SQL struct {
				Enabled string `json:"enabled"`
			} `json:"sql"`
		} `json:"xpack"`
		Rest struct {
			Action struct {
				Multi struct {
					AllowExplicitIndex string `json:"allow_explicit_index"`
				} `json:"multi"`
			} `json:"action"`
		} `json:"rest"`
		Cache struct {
			Recycler struct {
				Page struct {
					Limit struct {
						Heap string `json:"heap"`
					} `json:"limit"`
					Type   string `json:"type"`
					Weight struct {
						Longs   string `json:"longs"`
						Ints    string `json:"ints"`
						Bytes   string `json:"bytes"`
						Objects string `json:"objects"`
					} `json:"weight"`
				} `json:"page"`
			} `json:"recycler"`
		} `json:"cache"`
		Reindex struct {
			Remote struct {
				Whitelist []interface{} `json:"whitelist"`
			} `json:"remote"`
		} `json:"reindex"`
		Max struct {
			Anomaly struct {
				Records string `json:"records"`
			} `json:"anomaly"`
		} `json:"max"`
		Resource struct {
			Reload struct {
				Enabled  string `json:"enabled"`
				Interval struct {
					Low    string `json:"low"`
					High   string `json:"high"`
					Medium string `json:"medium"`
				} `json:"interval"`
			} `json:"reload"`
		} `json:"resource"`
		ThreadPool struct {
			ForceMerge struct {
				QueueSize string `json:"queue_size"`
				Size      string `json:"size"`
			} `json:"force_merge"`
			FetchShardStarted struct {
				Core      string `json:"core"`
				Max       string `json:"max"`
				KeepAlive string `json:"keep_alive"`
			} `json:"fetch_shard_started"`
			Listener struct {
				QueueSize string `json:"queue_size"`
				Size      string `json:"size"`
			} `json:"listener"`
			Index struct {
				QueueSize string `json:"queue_size"`
				Size      string `json:"size"`
			} `json:"index"`
			Refresh struct {
				Core      string `json:"core"`
				Max       string `json:"max"`
				KeepAlive string `json:"keep_alive"`
			} `json:"refresh"`
			Generic struct {
				Core      string `json:"core"`
				Max       string `json:"max"`
				KeepAlive string `json:"keep_alive"`
			} `json:"generic"`
			Warmer struct {
				Core      string `json:"core"`
				Max       string `json:"max"`
				KeepAlive string `json:"keep_alive"`
			} `json:"warmer"`
			Search struct {
				MaxQueueSize       string `json:"max_queue_size"`
				QueueSize          string `json:"queue_size"`
				Size               string `json:"size"`
				AutoQueueFrameSize string `json:"auto_queue_frame_size"`
				TargetResponseTime string `json:"target_response_time"`
				MinQueueSize       string `json:"min_queue_size"`
			} `json:"search"`
			FetchShardStore struct {
				Core      string `json:"core"`
				Max       string `json:"max"`
				KeepAlive string `json:"keep_alive"`
			} `json:"fetch_shard_store"`
			Flush struct {
				Core      string `json:"core"`
				Max       string `json:"max"`
				KeepAlive string `json:"keep_alive"`
			} `json:"flush"`
			Management struct {
				Core      string `json:"core"`
				Max       string `json:"max"`
				KeepAlive string `json:"keep_alive"`
			} `json:"management"`
			Analyze struct {
				QueueSize string `json:"queue_size"`
				Size      string `json:"size"`
			} `json:"analyze"`
			Get struct {
				QueueSize string `json:"queue_size"`
				Size      string `json:"size"`
			} `json:"get"`
			Bulk struct {
				QueueSize string `json:"queue_size"`
				Size      string `json:"size"`
			} `json:"bulk"`
			EstimatedTimeInterval string `json:"estimated_time_interval"`
			Write                 struct {
				QueueSize string `json:"queue_size"`
				Size      string `json:"size"`
			} `json:"write"`
			Snapshot struct {
				Core      string `json:"core"`
				Max       string `json:"max"`
				KeepAlive string `json:"keep_alive"`
			} `json:"snapshot"`
			SearchThrottled struct {
				MaxQueueSize       string `json:"max_queue_size"`
				QueueSize          string `json:"queue_size"`
				Size               string `json:"size"`
				AutoQueueFrameSize string `json:"auto_queue_frame_size"`
				TargetResponseTime string `json:"target_response_time"`
				MinQueueSize       string `json:"min_queue_size"`
			} `json:"search_throttled"`
		} `json:"thread_pool"`
		Index struct {
			Codec string `json:"codec"`
			Store struct {
				Type string `json:"type"`
				Fs   struct {
					FsLock string `json:"fs_lock"`
				} `json:"fs"`
				Preload []interface{} `json:"preload"`
			} `json:"store"`
		} `json:"index"`
		Monitor struct {
			Jvm struct {
				Gc struct {
					Enabled  string `json:"enabled"`
					Overhead struct {
						Warn  string `json:"warn"`
						Debug string `json:"debug"`
						Info  string `json:"info"`
					} `json:"overhead"`
					RefreshInterval string `json:"refresh_interval"`
				} `json:"gc"`
				RefreshInterval string `json:"refresh_interval"`
			} `json:"jvm"`
			Process struct {
				RefreshInterval string `json:"refresh_interval"`
			} `json:"process"`
			Os struct {
				RefreshInterval string `json:"refresh_interval"`
			} `json:"os"`
			Fs struct {
				RefreshInterval string `json:"refresh_interval"`
			} `json:"fs"`
		} `json:"monitor"`
		Transport struct {
			TCP struct {
				ReuseAddress      string `json:"reuse_address"`
				ConnectTimeout    string `json:"connect_timeout"`
				Compress          string `json:"compress"`
				Port              string `json:"port"`
				KeepAlive         string `json:"keep_alive"`
				ReceiveBufferSize string `json:"receive_buffer_size"`
				SendBufferSize    string `json:"send_buffer_size"`
			} `json:"tcp"`
			BindHost           []interface{} `json:"bind_host"`
			PingSchedule       string        `json:"ping_schedule"`
			ConnectionsPerNode struct {
				Recovery string `json:"recovery"`
				State    string `json:"state"`
				Bulk     string `json:"bulk"`
				Reg      string `json:"reg"`
				Ping     string `json:"ping"`
			} `json:"connections_per_node"`
			Tracer struct {
				Include []interface{} `json:"include"`
				Exclude []string      `json:"exclude"`
			} `json:"tracer"`
			Type        string `json:"type"`
			TypeDefault string `json:"type.default"`
			Features    struct {
				XPack string `json:"x-pack"`
			} `json:"features"`
			Host        []interface{} `json:"host"`
			PublishPort string        `json:"publish_port"`
			TCPNoDelay  string        `json:"tcp_no_delay"`
			PublishHost []interface{} `json:"publish_host"`
			Netty       struct {
				ReceivePredictorSize string `json:"receive_predictor_size"`
				ReceivePredictorMax  string `json:"receive_predictor_max"`
				WorkerCount          string `json:"worker_count"`
				ReceivePredictorMin  string `json:"receive_predictor_min"`
				BossCount            string `json:"boss_count"`
			} `json:"netty"`
		} `json:"transport"`
		Script struct {
			AllowedContexts     []interface{} `json:"allowed_contexts"`
			MaxCompilationsRate string        `json:"max_compilations_rate"`
			Cache               struct {
				MaxSize string `json:"max_size"`
				Expire  string `json:"expire"`
			} `json:"cache"`
			Painless struct {
				Regex struct {
					Enabled string `json:"enabled"`
				} `json:"regex"`
			} `json:"painless"`
			MaxSizeInBytes string        `json:"max_size_in_bytes"`
			AllowedTypes   []interface{} `json:"allowed_types"`
		} `json:"script"`
		Node struct {
			Data                          string `json:"data"`
			EnableLuceneSegmentInfosTrace string `json:"enable_lucene_segment_infos_trace"`
			LocalStorage                  string `json:"local_storage"`
			MaxLocalStorageNodes          string `json:"max_local_storage_nodes"`
			Name                          string `json:"name"`
			ID                            struct {
				Seed string `json:"seed"`
			} `json:"id"`
			Store struct {
				AllowMmapfs string `json:"allow_mmapfs"`
			} `json:"store"`
			Attr struct {
				Xpack struct {
					Installed string `json:"installed"`
				} `json:"xpack"`
				Ml struct {
					MachineMemory string `json:"machine_memory"`
					MaxOpenJobs   string `json:"max_open_jobs"`
					Enabled       string `json:"enabled"`
				} `json:"ml"`
			} `json:"attr"`
			Portsfile string `json:"portsfile"`
			Ingest    string `json:"ingest"`
			Master    string `json:"master"`
			Ml        string `json:"ml"`
		} `json:"node"`
		Indices struct {
			Cache struct {
				CleanupInterval string `json:"cleanup_interval"`
			} `json:"cache"`
			Mapping struct {
				DynamicTimeout string `json:"dynamic_timeout"`
			} `json:"mapping"`
			Memory struct {
				Interval           string `json:"interval"`
				MaxIndexBufferSize string `json:"max_index_buffer_size"`
				ShardInactiveTime  string `json:"shard_inactive_time"`
				IndexBufferSize    string `json:"index_buffer_size"`
				MinIndexBufferSize string `json:"min_index_buffer_size"`
			} `json:"memory"`
			Breaker struct {
				Request struct {
					Limit    string `json:"limit"`
					Type     string `json:"type"`
					Overhead string `json:"overhead"`
				} `json:"request"`
				Total struct {
					Limit string `json:"limit"`
				} `json:"total"`
				Accounting struct {
					Limit    string `json:"limit"`
					Overhead string `json:"overhead"`
				} `json:"accounting"`
				Fielddata struct {
					Limit    string `json:"limit"`
					Type     string `json:"type"`
					Overhead string `json:"overhead"`
				} `json:"fielddata"`
				Type string `json:"type"`
			} `json:"breaker"`
			Fielddata struct {
				Cache struct {
					Size string `json:"size"`
				} `json:"cache"`
			} `json:"fielddata"`
			Query struct {
				Bool struct {
					MaxClauseCount string `json:"max_clause_count"`
				} `json:"bool"`
				QueryString struct {
					AnalyzeWildcard      string `json:"analyze_wildcard"`
					AllowLeadingWildcard string `json:"allowLeadingWildcard"`
				} `json:"query_string"`
			} `json:"query"`
			Admin struct {
				FilteredFields string `json:"filtered_fields"`
			} `json:"admin"`
			Recovery struct {
				RecoveryActivityTimeout   string `json:"recovery_activity_timeout"`
				RetryDelayNetwork         string `json:"retry_delay_network"`
				InternalActionTimeout     string `json:"internal_action_timeout"`
				RetryDelayStateSync       string `json:"retry_delay_state_sync"`
				InternalActionLongTimeout string `json:"internal_action_long_timeout"`
				MaxBytesPerSec            string `json:"max_bytes_per_sec"`
			} `json:"recovery"`
			Requests struct {
				Cache struct {
					Size   string `json:"size"`
					Expire string `json:"expire"`
				} `json:"cache"`
			} `json:"requests"`
			Store struct {
				Delete struct {
					Shard struct {
						Timeout string `json:"timeout"`
					} `json:"shard"`
				} `json:"delete"`
			} `json:"store"`
			Analysis struct {
				Hunspell struct {
					Dictionary struct {
						IgnoreCase string `json:"ignore_case"`
						Lazy       string `json:"lazy"`
					} `json:"dictionary"`
				} `json:"hunspell"`
			} `json:"analysis"`
			Queries struct {
				Cache struct {
					Count       string `json:"count"`
					Size        string `json:"size"`
					AllSegments string `json:"all_segments"`
				} `json:"cache"`
			} `json:"queries"`
		} `json:"indices"`
		Plugin struct {
			Mandatory []interface{} `json:"mandatory"`
		} `json:"plugin"`
		MaxRunningJobs string `json:"max_running_jobs"`
		Discovery      struct {
			Type string `json:"type"`
			Zen  struct {
				CommitTimeout     string `json:"commit_timeout"`
				NoMasterBlock     string `json:"no_master_block"`
				JoinRetryDelay    string `json:"join_retry_delay"`
				JoinRetryAttempts string `json:"join_retry_attempts"`
				Ping              struct {
					Unicast struct {
						ConcurrentConnects  string        `json:"concurrent_connects"`
						Hosts               []interface{} `json:"hosts"`
						HostsResolveTimeout string        `json:"hosts.resolve_timeout"`
					} `json:"unicast"`
				} `json:"ping"`
				MasterElection struct {
					IgnoreNonMasterPings string `json:"ignore_non_master_pings"`
					WaitForJoinsTimeout  string `json:"wait_for_joins_timeout"`
				} `json:"master_election"`
				SendLeaveRequest string `json:"send_leave_request"`
				PingTimeout      string `json:"ping_timeout"`
				JoinTimeout      string `json:"join_timeout"`
				PublishDiff      struct {
					Enable string `json:"enable"`
				} `json:"publish_diff"`
				Publish struct {
					MaxPendingClusterStates string `json:"max_pending_cluster_states"`
				} `json:"publish"`
				MinimumMasterNodes string        `json:"minimum_master_nodes"`
				HostsProvider      []interface{} `json:"hosts_provider"`
				PublishTimeout     string        `json:"publish_timeout"`
				Fd                 struct {
					ConnectOnNetworkDisconnect string `json:"connect_on_network_disconnect"`
					PingInterval               string `json:"ping_interval"`
					PingRetries                string `json:"ping_retries"`
					RegisterConnectionListener string `json:"register_connection_listener"`
					PingTimeout                string `json:"ping_timeout"`
				} `json:"fd"`
				MaxPingsFromAnotherMaster string `json:"max_pings_from_another_master"`
			} `json:"zen"`
			InitialStateTimeout string `json:"initial_state_timeout"`
		} `json:"discovery"`
		Tribe struct {
			Name       string `json:"name"`
			OnConflict string `json:"on_conflict"`
			Blocks     struct {
				Metadata string `json:"metadata"`
				Read     struct {
					Indices []interface{} `json:"indices"`
				} `json:"read"`
				WriteIndices    []interface{} `json:"write.indices"`
				Write           string        `json:"write"`
				MetadataIndices []interface{} `json:"metadata.indices"`
			} `json:"blocks"`
		} `json:"tribe"`
		HTTP struct {
			Cors struct {
				MaxAge           string `json:"max-age"`
				AllowOrigin      string `json:"allow-origin"`
				AllowHeaders     string `json:"allow-headers"`
				AllowCredentials string `json:"allow-credentials"`
				AllowMethods     string `json:"allow-methods"`
				Enabled          string `json:"enabled"`
			} `json:"cors"`
			MaxChunkSize         string `json:"max_chunk_size"`
			CompressionLevel     string `json:"compression_level"`
			MaxInitialLineLength string `json:"max_initial_line_length"`
			Type                 string `json:"type"`
			Pipelining           string `json:"pipelining"`
			Enabled              string `json:"enabled"`
			TypeDefault          string `json:"type.default"`
			ContentType          struct {
				Required string `json:"required"`
			} `json:"content_type"`
			Host             []interface{} `json:"host"`
			PublishPort      string        `json:"publish_port"`
			ReadTimeout      string        `json:"read_timeout"`
			MaxContentLength string        `json:"max_content_length"`
			Netty            struct {
				ReceivePredictorSize         string `json:"receive_predictor_size"`
				MaxCompositeBufferComponents string `json:"max_composite_buffer_components"`
				ReceivePredictorMax          string `json:"receive_predictor_max"`
				WorkerCount                  string `json:"worker_count"`
				ReceivePredictorMin          string `json:"receive_predictor_min"`
			} `json:"netty"`
			TCP struct {
				ReuseAddress      string `json:"reuse_address"`
				KeepAlive         string `json:"keep_alive"`
				ReceiveBufferSize string `json:"receive_buffer_size"`
				SendBufferSize    string `json:"send_buffer_size"`
			} `json:"tcp"`
			BindHost              []interface{} `json:"bind_host"`
			ResetCookies          string        `json:"reset_cookies"`
			MaxWarningHeaderCount string        `json:"max_warning_header_count"`
			MaxWarningHeaderSize  string        `json:"max_warning_header_size"`
			DetailedErrors        struct {
				Enabled string `json:"enabled"`
			} `json:"detailed_errors"`
			Port                string        `json:"port"`
			MaxHeaderSize       string        `json:"max_header_size"`
			PipeliningMaxEvents string        `json:"pipelining.max_events"`
			TCPNoDelay          string        `json:"tcp_no_delay"`
			Compression         string        `json:"compression"`
			PublishHost         []interface{} `json:"publish_host"`
		} `json:"http"`
		Gateway struct {
			RecoverAfterMasterNodes string `json:"recover_after_master_nodes"`
			ExpectedNodes           string `json:"expected_nodes"`
			RecoverAfterDataNodes   string `json:"recover_after_data_nodes"`
			ExpectedDataNodes       string `json:"expected_data_nodes"`
			RecoverAfterTime        string `json:"recover_after_time"`
			ExpectedMasterNodes     string `json:"expected_master_nodes"`
			RecoverAfterNodes       string `json:"recover_after_nodes"`
		} `json:"gateway"`
	} `json:"defaults"`
}

type Settings struct {
	setting Setting
}

//节点设置
func NewSettings(client *elastic.Client) (*Settings, error) {
	settings := &Settings{}

	res, err := client.PerformRequest(context.Background(), elastic.PerformRequestOptions{
		Method: "GET",
		Path:   "/_cluster/settings?include_defaults=true",
	})
	if err != nil {
		return nil, err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(res.Body, &settings.setting)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// 获取配置文件设置的 PathRepo
func (this *Settings) GetPathRepo() []string {
	return this.setting.Defaults.Path.Repo
}

// 获取配置文件中的 AllowedUrls
func (this *Settings) GetAllowedUrls() []interface{} {
	return this.setting.Defaults.Repositories.URL.AllowedUrls
}
