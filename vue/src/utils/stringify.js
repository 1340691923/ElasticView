var SqlParser
!(function(e, s) {
  typeof exports === 'object' && typeof module !== 'undefined' ? s(exports) : typeof define === 'function' && define.amd ? define(['exports'], s) : s((e = e || self).SqlParser = {})
}(this, function(e) {
  'use strict'
  var s = typeof globalThis !== 'undefined' ? globalThis : typeof window !== 'undefined' ? window : typeof global !== 'undefined' ? global : typeof self !== 'undefined' ? self : {}

  function n(e, s) {
    return e(s = { exports: {}}, s.exports), s.exports
  }

  var t = n(function(e) {
    var n, t
    n = s, t = function() {
      function e(s, n, t) {
        return this.id = ++e.highestId, this.name = s, this.symbols = n, this.postprocess = t, this
      }

      function s(e, s, n, t) {
        this.rule = e, this.dot = s, this.reference = n, this.data = [], this.wantedBy = t, this.isComplete = this.dot === e.symbols.length
      }

      function n(e, s) {
        this.grammar = e, this.index = s, this.states = [], this.wants = {}, this.scannable = [], this.completed = {}
      }

      function t(e, s) {
        this.rules = e, this.start = s || this.rules[0].name
        var n = this.byName = {}
        this.rules.forEach(function(e) {
          n.hasOwnProperty(e.name) || (n[e.name] = []), n[e.name].push(e)
        })
      }

      function o() {
        this.reset('')
      }

      function r(e, s, r) {
        if (e instanceof t) {
          var i = e
          r = s
        } else {
          i = t.fromCompiled(e, s)
        }
        for (var p in this.grammar = i, this.options = {
          keepHistory: !1,
          lexer: i.lexer || new o()
        }, r || {}) {
          this.options[p] = r[p]
        }
        this.lexer = this.options.lexer, this.lexerState = void 0
        var $ = new n(i, 0)
        this.table = [$], $.wants[i.start] = [], $.predict(i.start), $.process(), this.current = 0
      }

      function i(e) {
        var s = typeof e
        if (s === 'string') return e
        if (s === 'object') {
          if (e.literal) return JSON.stringify(e.literal)
          if (e instanceof RegExp) return e.toString()
          if (e.type) return '%' + e.type
          if (e.test) return '<' + String(e.test) + '>'
          throw new Error('Unknown symbol type: ' + e)
        }
      }

      return e.highestId = 0, e.prototype.toString = function(e) {
        var s = void 0 === e ? this.symbols.map(i).join(' ') : this.symbols.slice(0, e).map(i).join(' ') + ' ● ' + this.symbols.slice(e).map(i).join(' ')
        return this.name + ' → ' + s
      }, s.prototype.toString = function() {
        return '{' + this.rule.toString(this.dot) + '}, from: ' + (this.reference || 0)
      }, s.prototype.nextState = function(e) {
        var n = new s(this.rule, this.dot + 1, this.reference, this.wantedBy)
        return n.left = this, n.right = e, n.isComplete && (n.data = n.build(), n.right = void 0), n
      }, s.prototype.build = function() {
        var e = []; var s = this
        do {
          e.push(s.right.data), s = s.left
        } while (s.left)
        return e.reverse(), e
      }, s.prototype.finish = function() {
        this.rule.postprocess && (this.data = this.rule.postprocess(this.data, this.reference, r.fail))
      }, n.prototype.process = function(e) {
        for (var s = this.states, n = this.wants, t = this.completed, o = 0; o < s.length; o++) {
          var i = s[o]
          if (i.isComplete) {
            if (i.finish(), i.data !== r.fail) {
              for (var p = i.wantedBy, $ = p.length; $--;) {
                var a = p[$]
                this.complete(a, i)
              }
              if (i.reference === this.index) {
                var l = i.rule.name;
                (this.completed[l] = this.completed[l] || []).push(i)
              }
            }
          } else {
            if (typeof (l = i.rule.symbols[i.dot]) !== 'string') {
              this.scannable.push(i)
              continue
            }
            if (n[l]) {
              if (n[l].push(i), t.hasOwnProperty(l)) {
                var b = t[l]
                for ($ = 0; $ < b.length; $++) {
                  var u = b[$]
                  this.complete(i, u)
                }
              }
            } else {
              n[l] = [i], this.predict(l)
            }
          }
        }
      }, n.prototype.predict = function(e) {
        for (var n = this.grammar.byName[e] || [], t = 0; t < n.length; t++) {
          var o = n[t]; var r = this.wants[e]; var i = new s(o, 0, this.index, r)
          this.states.push(i)
        }
      }, n.prototype.complete = function(e, s) {
        var n = e.nextState(s)
        this.states.push(n)
      }, t.fromCompiled = function(s, n) {
        var o = s.Lexer
        s.ParserStart && (n = s.ParserStart, s = s.ParserRules)
        var r = new t(s = s.map(function(s) {
          return new e(s.name, s.symbols, s.postprocess)
        }), n)
        return r.lexer = o, r
      }, o.prototype.reset = function(e, s) {
        this.buffer = e, this.index = 0, this.line = s ? s.line : 1, this.lastLineBreak = s ? -s.col : 0
      }, o.prototype.next = function() {
        if (this.index < this.buffer.length) {
          var e = this.buffer[this.index++]
          return e === '\n' && (this.line += 1, this.lastLineBreak = this.index), { value: e }
        }
      }, o.prototype.save = function() {
        return { line: this.line, col: this.index - this.lastLineBreak }
      }, o.prototype.formatError = function(e, s) {
        var n = this.buffer
        if (typeof n === 'string') {
          var t = n.split('\n').slice(Math.max(0, this.line - 5), this.line); var o = n.indexOf('\n', this.index)
          o === -1 && (o = n.length)
          var r = this.index - this.lastLineBreak; var i = String(this.line).length
          return s += ' at line ' + this.line + ' col ' + r + ':\n\n', s += t.map(function(e, s) {
            return p(this.line - t.length + s + 1, i) + ' ' + e
          }, this).join('\n'), s += '\n' + p('', i + r) + '^\n'
        }
        return s + ' at index ' + (this.index - 1)

        function p(e, s) {
          var n = String(e)
          return Array(s - n.length + 1).join(' ') + n
        }
      }, r.fail = {}, r.prototype.feed = function(e) {
        var s; var t = this.lexer
        for (t.reset(e, this.lexerState); ;) {
          try {
            if (!(s = t.next())) break
          } catch (e) {
            var r = new n(this.grammar, this.current + 1)
            throw this.table.push(r), ($ = new Error(this.reportLexerError(e))).offset = this.current, $.token = e.token, $
          }
          var i = this.table[this.current]
          this.options.keepHistory || delete this.table[this.current - 1]
          var p = this.current + 1
          r = new n(this.grammar, p), this.table.push(r)
          for (var $, a = void 0 !== s.text ? s.text : s.value, l = t.constructor === o ? s.value : s, b = i.scannable, u = b.length; u--;) {
            var m = b[u]; var c = m.rule.symbols[m.dot]
            if (c.test ? c.test(l) : c.type ? c.type === s.type : c.literal === a) {
              var _ = m.nextState({ data: l, token: s, isToken: !0, reference: p - 1 })
              r.states.push(_)
            }
          }
          if (r.process(), r.states.length === 0) throw ($ = new Error(this.reportError(s))).offset = this.current, $.token = s, $
          this.options.keepHistory && (i.lexerState = t.save()), this.current++
        }
        return i && (this.lexerState = t.save()), this.results = this.finish(), this
      }, r.prototype.reportLexerError = function(e) {
        var s; var n; var t = e.token
        return t ? (s = 'input ' + JSON.stringify(t.text[0]) + ' (lexer error)', n = this.lexer.formatError(t, 'Syntax error')) : (s = 'input (lexer error)', n = e.message), this.reportErrorCommon(n, s)
      }, r.prototype.reportError = function(e) {
        var s = (e.type ? e.type + ' token: ' : '') + JSON.stringify(void 0 !== e.value ? e.value : e)
        var n = this.lexer.formatError(e, 'Syntax error')
        return this.reportErrorCommon(n, s)
      }, r.prototype.reportErrorCommon = function(e, s) {
        var n = []
        n.push(e)
        var t = this.table.length - 2; var o = this.table[t]; var r = o.states.filter(function(e) {
          var s = e.rule.symbols[e.dot]
          return s && typeof s !== 'string'
        })
        return r.length === 0 ? (n.push('Unexpected ' + s + '. I did not expect any more input. Here is the state of my parse table:\n'), this.displayStateStack(o.states, n)) : (n.push('Unexpected ' + s + '. Instead, I was expecting to see one of the following:\n'), r.map(function(e) {
          return this.buildFirstStateStack(e, []) || [e]
        }, this).forEach(function(e) {
          var s = e[0]; var t = s.rule.symbols[s.dot]; var o = this.getSymbolDisplay(t)
          n.push('A ' + o + ' based on:'), this.displayStateStack(e, n)
        }, this)), n.push(''), n.join('\n')
      }, r.prototype.displayStateStack = function(e, s) {
        for (var n, t = 0, o = 0; o < e.length; o++) {
          var r = e[o]; var i = r.rule.toString(r.dot)
          i === n ? t++ : (t > 0 && s.push('    ^ ' + t + ' more lines identical to this'), t = 0, s.push('    ' + i)), n = i
        }
      }, r.prototype.getSymbolDisplay = function(e) {
        return (function(e) {
          var s = typeof e
          if (s === 'string') return e
          if (s === 'object') {
            if (e.literal) return JSON.stringify(e.literal)
            if (e instanceof RegExp) return 'character matching ' + e
            if (e.type) return e.type + ' token'
            if (e.test) return 'token matching ' + String(e.test)
            throw new Error('Unknown symbol type: ' + e)
          }
        }(e))
      }, r.prototype.buildFirstStateStack = function(e, s) {
        if (s.indexOf(e) !== -1) return null
        if (e.wantedBy.length === 0) return [e]
        var n = e.wantedBy[0]; var t = [e].concat(s); var o = this.buildFirstStateStack(n, t)
        return o === null ? null : [e].concat(o)
      }, r.prototype.save = function() {
        var e = this.table[this.current]
        return e.lexerState = this.lexerState, e
      }, r.prototype.restore = function(e) {
        var s = e.index
        this.current = s, this.table[s] = e, this.table.splice(s + 1), this.lexerState = e.lexerState, this.results = this.finish()
      }, r.prototype.rewind = function(e) {
        if (!this.options.keepHistory) throw new Error('set option `keepHistory` to enable rewinding')
        this.restore(this.table[e])
      }, r.prototype.finish = function() {
        var e = []; var s = this.grammar.start
        return this.table[this.table.length - 1].states.forEach(function(n) {
          n.rule.name === s && n.dot === n.rule.symbols.length && n.reference === 0 && n.data !== r.fail && e.push(n)
        }), e.map(function(e) {
          return e.data
        })
      }, { Parser: r, Grammar: t, Rule: e }
    }, e.exports ? e.exports = t() : n.nearley = t()
  })
  var o = n(function(e) {
    !(function() {
      function s(e) {
        return e[0]
      }

      var n = {
        Lexer: void 0,
        ParserRules: [{ name: '_$ebnf$1', symbols: [] }, {
          name: '_$ebnf$1',
          symbols: ['_$ebnf$1', 'wschar'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: '_', symbols: ['_$ebnf$1'], postprocess: function(e) {
            return null
          }
        }, { name: '__$ebnf$1', symbols: ['wschar'] }, {
          name: '__$ebnf$1',
          symbols: ['__$ebnf$1', 'wschar'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: '__', symbols: ['__$ebnf$1'], postprocess: function(e) {
            return null
          }
        }, { name: 'wschar', symbols: [/[ \t\n\v\f]/], postprocess: s }, {
          name: 'unsigned_int$ebnf$1',
          symbols: [/[0-9]/]
        }, {
          name: 'unsigned_int$ebnf$1', symbols: ['unsigned_int$ebnf$1', /[0-9]/], postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'unsigned_int', symbols: ['unsigned_int$ebnf$1'], postprocess: function(e) {
            return parseInt(e[0].join(''))
          }
        }, { name: 'int$ebnf$1$subexpression$1', symbols: [{ literal: '-' }] }, {
          name: 'int$ebnf$1$subexpression$1',
          symbols: [{ literal: '+' }]
        }, { name: 'int$ebnf$1', symbols: ['int$ebnf$1$subexpression$1'], postprocess: s }, {
          name: 'int$ebnf$1',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, { name: 'int$ebnf$2', symbols: [/[0-9]/] }, {
          name: 'int$ebnf$2',
          symbols: ['int$ebnf$2', /[0-9]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'int', symbols: ['int$ebnf$1', 'int$ebnf$2'], postprocess: function(e) {
            return e[0] ? parseInt(e[0][0] + e[1].join('')) : parseInt(e[1].join(''))
          }
        }, { name: 'unsigned_decimal$ebnf$1', symbols: [/[0-9]/] }, {
          name: 'unsigned_decimal$ebnf$1',
          symbols: ['unsigned_decimal$ebnf$1', /[0-9]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'unsigned_decimal$ebnf$2$subexpression$1$ebnf$1',
          symbols: [/[0-9]/]
        }, {
          name: 'unsigned_decimal$ebnf$2$subexpression$1$ebnf$1',
          symbols: ['unsigned_decimal$ebnf$2$subexpression$1$ebnf$1', /[0-9]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'unsigned_decimal$ebnf$2$subexpression$1',
          symbols: [{ literal: '.' }, 'unsigned_decimal$ebnf$2$subexpression$1$ebnf$1']
        }, {
          name: 'unsigned_decimal$ebnf$2',
          symbols: ['unsigned_decimal$ebnf$2$subexpression$1'],
          postprocess: s
        }, {
          name: 'unsigned_decimal$ebnf$2', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'unsigned_decimal',
          symbols: ['unsigned_decimal$ebnf$1', 'unsigned_decimal$ebnf$2'],
          postprocess: function(e) {
            return parseFloat(e[0].join('') + (e[1] ? '.' + e[1][1].join('') : ''))
          }
        }, { name: 'decimal$ebnf$1', symbols: [{ literal: '-' }], postprocess: s }, {
          name: 'decimal$ebnf$1',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, { name: 'decimal$ebnf$2', symbols: [/[0-9]/] }, {
          name: 'decimal$ebnf$2',
          symbols: ['decimal$ebnf$2', /[0-9]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'decimal$ebnf$3$subexpression$1$ebnf$1',
          symbols: [/[0-9]/]
        }, {
          name: 'decimal$ebnf$3$subexpression$1$ebnf$1',
          symbols: ['decimal$ebnf$3$subexpression$1$ebnf$1', /[0-9]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'decimal$ebnf$3$subexpression$1',
          symbols: [{ literal: '.' }, 'decimal$ebnf$3$subexpression$1$ebnf$1']
        }, {
          name: 'decimal$ebnf$3',
          symbols: ['decimal$ebnf$3$subexpression$1'],
          postprocess: s
        }, {
          name: 'decimal$ebnf$3', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'decimal',
          symbols: ['decimal$ebnf$1', 'decimal$ebnf$2', 'decimal$ebnf$3'],
          postprocess: function(e) {
            return parseFloat((e[0] || '') + e[1].join('') + (e[2] ? '.' + e[2][1].join('') : ''))
          }
        }, {
          name: 'percentage', symbols: ['decimal', { literal: '%' }], postprocess: function(e) {
            return e[0] / 100
          }
        }, { name: 'jsonfloat$ebnf$1', symbols: [{ literal: '-' }], postprocess: s }, {
          name: 'jsonfloat$ebnf$1',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, { name: 'jsonfloat$ebnf$2', symbols: [/[0-9]/] }, {
          name: 'jsonfloat$ebnf$2',
          symbols: ['jsonfloat$ebnf$2', /[0-9]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'jsonfloat$ebnf$3$subexpression$1$ebnf$1',
          symbols: [/[0-9]/]
        }, {
          name: 'jsonfloat$ebnf$3$subexpression$1$ebnf$1',
          symbols: ['jsonfloat$ebnf$3$subexpression$1$ebnf$1', /[0-9]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'jsonfloat$ebnf$3$subexpression$1',
          symbols: [{ literal: '.' }, 'jsonfloat$ebnf$3$subexpression$1$ebnf$1']
        }, {
          name: 'jsonfloat$ebnf$3',
          symbols: ['jsonfloat$ebnf$3$subexpression$1'],
          postprocess: s
        }, {
          name: 'jsonfloat$ebnf$3', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'jsonfloat$ebnf$4$subexpression$1$ebnf$1',
          symbols: [/[+-]/],
          postprocess: s
        }, {
          name: 'jsonfloat$ebnf$4$subexpression$1$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'jsonfloat$ebnf$4$subexpression$1$ebnf$2',
          symbols: [/[0-9]/]
        }, {
          name: 'jsonfloat$ebnf$4$subexpression$1$ebnf$2',
          symbols: ['jsonfloat$ebnf$4$subexpression$1$ebnf$2', /[0-9]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'jsonfloat$ebnf$4$subexpression$1',
          symbols: [/[eE]/, 'jsonfloat$ebnf$4$subexpression$1$ebnf$1', 'jsonfloat$ebnf$4$subexpression$1$ebnf$2']
        }, {
          name: 'jsonfloat$ebnf$4',
          symbols: ['jsonfloat$ebnf$4$subexpression$1'],
          postprocess: s
        }, {
          name: 'jsonfloat$ebnf$4', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'jsonfloat',
          symbols: ['jsonfloat$ebnf$1', 'jsonfloat$ebnf$2', 'jsonfloat$ebnf$3', 'jsonfloat$ebnf$4'],
          postprocess: function(e) {
            return parseFloat((e[0] || '') + e[1].join('') + (e[2] ? '.' + e[2][1].join('') : '') + (e[3] ? 'e' + (e[3][1] || '+') + e[3][2].join('') : ''))
          }
        }, { name: 'dqstring$ebnf$1', symbols: [] }, {
          name: 'dqstring$ebnf$1',
          symbols: ['dqstring$ebnf$1', 'dstrchar'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'dqstring',
          symbols: [{ literal: '"' }, 'dqstring$ebnf$1', { literal: '"' }],
          postprocess: function(e) {
            return e[1].join('')
          }
        }, { name: 'sqstring$ebnf$1', symbols: [] }, {
          name: 'sqstring$ebnf$1',
          symbols: ['sqstring$ebnf$1', 'sstrchar'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'sqstring',
          symbols: [{ literal: '\'' }, 'sqstring$ebnf$1', { literal: '\'' }],
          postprocess: function(e) {
            return e[1].join('')
          }
        }, { name: 'btstring$ebnf$1', symbols: [] }, {
          name: 'btstring$ebnf$1',
          symbols: ['btstring$ebnf$1', /[^`]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'btstring',
          symbols: [{ literal: '`' }, 'btstring$ebnf$1', { literal: '`' }],
          postprocess: function(e) {
            return e[1].join('')
          }
        }, { name: 'dstrchar', symbols: [/[^\\"\n]/], postprocess: s }, {
          name: 'dstrchar',
          symbols: [{ literal: '\\' }, 'strescape'],
          postprocess: function(e) {
            return JSON.parse('"' + e.join('') + '"')
          }
        }, { name: 'sstrchar', symbols: [/[^\\'\n]/], postprocess: s }, {
          name: 'sstrchar',
          symbols: [{ literal: '\\' }, 'strescape'],
          postprocess: function(e) {
            return JSON.parse('"' + e.join('') + '"')
          }
        }, {
          name: 'sstrchar$string$1', symbols: [{ literal: '\\' }, { literal: '\'' }], postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'sstrchar', symbols: ['sstrchar$string$1'], postprocess: function(e) {
            return '\''
          }
        }, { name: 'strescape', symbols: [/["\\/bfnrt]/], postprocess: s }, {
          name: 'strescape',
          symbols: [{ literal: 'u' }, /[a-fA-F0-9]/, /[a-fA-F0-9]/, /[a-fA-F0-9]/, /[a-fA-F0-9]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, { name: 'table_name', symbols: ['word'], postprocess: e => e[0] }, {
          name: 'table_name',
          symbols: ['btstring'],
          postprocess: e => e[0]
        }, { name: 'field_name', symbols: ['word'], postprocess: e => e[0] }, {
          name: 'field_name',
          symbols: ['btstring'],
          postprocess: e => e[0]
        }, { name: 'set_value', symbols: ['word'], postprocess: e => e[0] }, {
          name: 'set_value',
          symbols: ['btstring'],
          postprocess: e => e[0]
        }, { name: 'string', symbols: ['dqstring'], postprocess: e => e[0] }, {
          name: 'string',
          symbols: ['sqstring'],
          postprocess: e => e[0]
        }, { name: 'word', symbols: [/[A-Za-z_]/], postprocess: s }, {
          name: 'word',
          symbols: ['word', /[A-Za-z0-9_]/],
          postprocess: e => '' + e[0] + e[1]
        }, { name: 'newline$ebnf$1', symbols: [] }, {
          name: 'newline$ebnf$1',
          symbols: ['newline$ebnf$1', /[ \t\v\f]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, { name: 'newline', symbols: ['newline$ebnf$1', /[\r?\n]/], postprocess: e => null }, {
          name: 'space$ebnf$1',
          symbols: []
        }, {
          name: 'space$ebnf$1', symbols: ['space$ebnf$1', /[ \t\v\f]/], postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, { name: 'space', symbols: ['space$ebnf$1'], postprocess: e => null }, {
          name: 'data_type',
          symbols: ['numeric_type'],
          postprocess: e => e[0]
        }, { name: 'data_type', symbols: ['datetime_type'], postprocess: e => e[0] }, {
          name: 'data_type',
          symbols: ['string_type'],
          postprocess: e => e[0]
        }, { name: 'data_type', symbols: ['spatial_data_type'], postprocess: e => e[0] }, {
          name: 'numeric_type',
          symbols: ['integer_type'],
          postprocess: e => e[0]
        }, {
          name: 'numeric_type$subexpression$1$subexpression$1',
          symbols: [/[dD]/, /[eE]/, /[cC]/, /[iI]/, /[mM]/, /[aA]/, /[lL]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'numeric_type$subexpression$1',
          symbols: ['numeric_type$subexpression$1$subexpression$1']
        }, {
          name: 'numeric_type$subexpression$1$subexpression$2',
          symbols: [/[nN]/, /[uU]/, /[mM]/, /[eE]/, /[rR]/, /[iI]/, /[cC]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'numeric_type$subexpression$1',
          symbols: ['numeric_type$subexpression$1$subexpression$2']
        }, {
          name: 'numeric_type$subexpression$1$subexpression$3',
          symbols: [/[fF]/, /[lL]/, /[oO]/, /[aA]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'numeric_type$subexpression$1',
          symbols: ['numeric_type$subexpression$1$subexpression$3']
        }, {
          name: 'numeric_type$subexpression$1$subexpression$4',
          symbols: [/[dD]/, /[oO]/, /[uU]/, /[bB]/, /[lL]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'numeric_type$subexpression$1',
          symbols: ['numeric_type$subexpression$1$subexpression$4']
        }, {
          name: 'numeric_type$ebnf$1$subexpression$1$subexpression$1',
          symbols: ['unsigned_int']
        }, {
          name: 'numeric_type$ebnf$1$subexpression$1$subexpression$1',
          symbols: ['unsigned_int', '_', { literal: ',' }, '_', 'unsigned_int']
        }, {
          name: 'numeric_type$ebnf$1$subexpression$1',
          symbols: ['_', { literal: '(' }, '_', 'numeric_type$ebnf$1$subexpression$1$subexpression$1', '_', { literal: ')' }, '_']
        }, {
          name: 'numeric_type$ebnf$1',
          symbols: ['numeric_type$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'numeric_type$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'numeric_type', symbols: ['numeric_type$subexpression$1', 'numeric_type$ebnf$1'], postprocess: e => {
            let s = e[0][0].toUpperCase()
            s === 'NUMERIC' && (s = 'DECIMAL')
            let n = { DECIMAL: 10, FLOAT: 23, DOUBLE: 53 }[s]; let t = 0
            return e[1] && e[1][3] && (e[1][3].length === 1 ? n = e[1][3][0] : e[1][3].length === 5 && (n = e[1][3][0], t = e[1][3][4])), {
              type: s,
              params: [n, t]
            }
          }
        }, {
          name: 'numeric_type$subexpression$2', symbols: [/[bB]/, /[iI]/, /[tT]/], postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'numeric_type',
          symbols: ['numeric_type$subexpression$2', '_', { literal: '(' }, '_', 'unsigned_int', '_', { literal: ')' }],
          postprocess: e => ({ type: 'BIT', params: [e[4]] })
        }, {
          name: 'integer_type$subexpression$1$subexpression$1',
          symbols: [/[tT]/, /[iI]/, /[nN]/, /[yY]/, /[iI]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'integer_type$subexpression$1',
          symbols: ['integer_type$subexpression$1$subexpression$1']
        }, {
          name: 'integer_type$subexpression$1$subexpression$2',
          symbols: [/[sS]/, /[mM]/, /[aA]/, /[lL]/, /[lL]/, /[iI]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'integer_type$subexpression$1',
          symbols: ['integer_type$subexpression$1$subexpression$2']
        }, {
          name: 'integer_type$subexpression$1$subexpression$3',
          symbols: [/[mM]/, /[eE]/, /[dD]/, /[iI]/, /[uU]/, /[mM]/, /[iI]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'integer_type$subexpression$1',
          symbols: ['integer_type$subexpression$1$subexpression$3']
        }, {
          name: 'integer_type$subexpression$1$subexpression$4',
          symbols: [/[iI]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'integer_type$subexpression$1',
          symbols: ['integer_type$subexpression$1$subexpression$4']
        }, {
          name: 'integer_type$subexpression$1$subexpression$5',
          symbols: [/[bB]/, /[iI]/, /[gG]/, /[iI]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'integer_type$subexpression$1',
          symbols: ['integer_type$subexpression$1$subexpression$5']
        }, {
          name: 'integer_type$ebnf$1$subexpression$1',
          symbols: ['_', { literal: '(' }, '_', 'unsigned_int', '_', { literal: ')' }, '_']
        }, {
          name: 'integer_type$ebnf$1',
          symbols: ['integer_type$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'integer_type$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'integer_type$ebnf$2$subexpression$1$subexpression$1',
          symbols: [/[uU]/, /[nN]/, /[sS]/, /[iI]/, /[gG]/, /[nN]/, /[eE]/, /[dD]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'integer_type$ebnf$2$subexpression$1',
          symbols: ['__', 'integer_type$ebnf$2$subexpression$1$subexpression$1']
        }, {
          name: 'integer_type$ebnf$2',
          symbols: ['integer_type$ebnf$2$subexpression$1'],
          postprocess: s
        }, {
          name: 'integer_type$ebnf$2', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'integer_type$ebnf$3$subexpression$1$subexpression$1',
          symbols: [/[zZ]/, /[eE]/, /[rR]/, /[oO]/, /[fF]/, /[iI]/, /[lL]/, /[lL]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'integer_type$ebnf$3$subexpression$1',
          symbols: ['__', 'integer_type$ebnf$3$subexpression$1$subexpression$1']
        }, {
          name: 'integer_type$ebnf$3',
          symbols: ['integer_type$ebnf$3$subexpression$1'],
          postprocess: s
        }, {
          name: 'integer_type$ebnf$3', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'integer_type',
          symbols: ['integer_type$subexpression$1', 'integer_type$ebnf$1', 'integer_type$ebnf$2', 'integer_type$ebnf$3'],
          postprocess: e => {
            const s = e[0][0].toUpperCase(); const n = {
              type: s,
              params: [e[1] ? e[1][3] : { TINYINT: 4, SMALLINT: 6, MEDIUMINT: 8, INT: 11, BIGINT: 20 }[s]]
            }
            return e[2] && e[2][1] && (n.unsigned = !0), e[3] && e[3][1] && (n.zerofill = !0), n
          }
        }, {
          name: 'datetime_type$subexpression$1',
          symbols: [/[dD]/, /[aA]/, /[tT]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'datetime_type',
          symbols: ['datetime_type$subexpression$1'],
          postprocess: e => ({ type: 'DATE' })
        }, {
          name: 'datetime_type$subexpression$2$subexpression$1',
          symbols: [/[dD]/, /[aA]/, /[tT]/, /[eE]/, /[tT]/, /[iI]/, /[mM]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'datetime_type$subexpression$2',
          symbols: ['datetime_type$subexpression$2$subexpression$1']
        }, {
          name: 'datetime_type$subexpression$2$subexpression$2',
          symbols: [/[tT]/, /[iI]/, /[mM]/, /[eE]/, /[sS]/, /[tT]/, /[aA]/, /[mM]/, /[pP]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'datetime_type$subexpression$2',
          symbols: ['datetime_type$subexpression$2$subexpression$2']
        }, {
          name: 'datetime_type$subexpression$2$subexpression$3',
          symbols: [/[tT]/, /[iI]/, /[mM]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'datetime_type$subexpression$2',
          symbols: ['datetime_type$subexpression$2$subexpression$3']
        }, {
          name: 'datetime_type$subexpression$2$subexpression$4',
          symbols: [/[yY]/, /[eE]/, /[aA]/, /[rR]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'datetime_type$subexpression$2',
          symbols: ['datetime_type$subexpression$2$subexpression$4']
        }, {
          name: 'datetime_type$ebnf$1$subexpression$1',
          symbols: ['_', { literal: '(' }, '_', 'unsigned_int', '_', { literal: ')' }, '_']
        }, {
          name: 'datetime_type$ebnf$1',
          symbols: ['datetime_type$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'datetime_type$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'datetime_type',
          symbols: ['datetime_type$subexpression$2', 'datetime_type$ebnf$1'],
          postprocess: e => {
            const s = []
            e[1] && e[1][3] && s.push(e[1][3])
            const n = { type: e[0][0].toUpperCase() }
            return s && (n.params = s), n
          }
        }, {
          name: 'string_type$subexpression$1$subexpression$1',
          symbols: [/[cC]/, /[hH]/, /[aA]/, /[rR]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$1',
          symbols: ['string_type$subexpression$1$subexpression$1']
        }, {
          name: 'string_type$subexpression$1$subexpression$2',
          symbols: [/[vV]/, /[aA]/, /[rR]/, /[cC]/, /[hH]/, /[aA]/, /[rR]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$1',
          symbols: ['string_type$subexpression$1$subexpression$2']
        }, {
          name: 'string_type$subexpression$1$subexpression$3',
          symbols: [/[bB]/, /[iI]/, /[nN]/, /[aA]/, /[rR]/, /[yY]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$1',
          symbols: ['string_type$subexpression$1$subexpression$3']
        }, {
          name: 'string_type$subexpression$1$subexpression$4',
          symbols: [/[vV]/, /[aA]/, /[rR]/, /[bB]/, /[iI]/, /[nN]/, /[aA]/, /[rR]/, /[yY]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$1',
          symbols: ['string_type$subexpression$1$subexpression$4']
        }, {
          name: 'string_type',
          symbols: ['string_type$subexpression$1', '_', { literal: '(' }, '_', 'unsigned_int', '_', { literal: ')' }],
          postprocess: e => ({ type: e[0][0].toUpperCase(), params: [e[4]] })
        }, {
          name: 'string_type$subexpression$2$subexpression$1',
          symbols: [/[tT]/, /[iI]/, /[nN]/, /[yY]/, /[bB]/, /[lL]/, /[oO]/, /[bB]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$1']
        }, {
          name: 'string_type$subexpression$2$subexpression$2',
          symbols: [/[bB]/, /[lL]/, /[oO]/, /[bB]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$2']
        }, {
          name: 'string_type$subexpression$2$subexpression$3',
          symbols: [/[mM]/, /[eE]/, /[dD]/, /[iI]/, /[uU]/, /[mM]/, /[bB]/, /[lL]/, /[oO]/, /[bB]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$3']
        }, {
          name: 'string_type$subexpression$2$subexpression$4',
          symbols: [/[lL]/, /[oO]/, /[nN]/, /[gG]/, /[bB]/, /[lL]/, /[oO]/, /[bB]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$4']
        }, {
          name: 'string_type$subexpression$2$subexpression$5',
          symbols: [/[tT]/, /[iI]/, /[nN]/, /[yY]/, /[tT]/, /[eE]/, /[xX]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$5']
        }, {
          name: 'string_type$subexpression$2$subexpression$6',
          symbols: [/[tT]/, /[eE]/, /[xX]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$6']
        }, {
          name: 'string_type$subexpression$2$subexpression$7',
          symbols: [/[mM]/, /[eE]/, /[dD]/, /[iI]/, /[uU]/, /[mM]/, /[tT]/, /[eE]/, /[xX]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$7']
        }, {
          name: 'string_type$subexpression$2$subexpression$8',
          symbols: [/[lL]/, /[oO]/, /[nN]/, /[gG]/, /[tT]/, /[eE]/, /[xX]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$8']
        }, {
          name: 'string_type$subexpression$2$subexpression$9',
          symbols: [/[jJ]/, /[sS]/, /[oO]/, /[nN]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$2',
          symbols: ['string_type$subexpression$2$subexpression$9']
        }, {
          name: 'string_type',
          symbols: ['string_type$subexpression$2'],
          postprocess: e => ({ type: e[0][0].toUpperCase() })
        }, {
          name: 'string_type$subexpression$3$subexpression$1',
          symbols: [/[eE]/, /[nN]/, /[uU]/, /[mM]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$3',
          symbols: ['string_type$subexpression$3$subexpression$1']
        }, {
          name: 'string_type$subexpression$3$subexpression$2',
          symbols: [/[sS]/, /[eE]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'string_type$subexpression$3',
          symbols: ['string_type$subexpression$3$subexpression$2']
        }, {
          name: 'string_type$ebnf$1$subexpression$1$ebnf$1$subexpression$1',
          symbols: [{ literal: ',' }, '_']
        }, {
          name: 'string_type$ebnf$1$subexpression$1$ebnf$1',
          symbols: ['string_type$ebnf$1$subexpression$1$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'string_type$ebnf$1$subexpression$1$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'string_type$ebnf$1$subexpression$1',
          symbols: ['string', '_', 'string_type$ebnf$1$subexpression$1$ebnf$1']
        }, {
          name: 'string_type$ebnf$1',
          symbols: ['string_type$ebnf$1$subexpression$1']
        }, {
          name: 'string_type$ebnf$1$subexpression$2$ebnf$1$subexpression$1',
          symbols: [{ literal: ',' }, '_']
        }, {
          name: 'string_type$ebnf$1$subexpression$2$ebnf$1',
          symbols: ['string_type$ebnf$1$subexpression$2$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'string_type$ebnf$1$subexpression$2$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'string_type$ebnf$1$subexpression$2',
          symbols: ['string', '_', 'string_type$ebnf$1$subexpression$2$ebnf$1']
        }, {
          name: 'string_type$ebnf$1',
          symbols: ['string_type$ebnf$1', 'string_type$ebnf$1$subexpression$2'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'string_type',
          symbols: ['string_type$subexpression$3', '_', { literal: '(' }, '_', 'string_type$ebnf$1', '_', { literal: ')' }],
          postprocess: e => ({ type: e[0][0].toUpperCase(), params: e[4].map(e => e[0]) })
        }, {
          name: 'spatial_data_type$subexpression$1$subexpression$1',
          symbols: [/[gG]/, /[eE]/, /[oO]/, /[mM]/, /[eE]/, /[tT]/, /[rR]/, /[yY]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'spatial_data_type$subexpression$1',
          symbols: ['spatial_data_type$subexpression$1$subexpression$1']
        }, {
          name: 'spatial_data_type$subexpression$1$subexpression$2',
          symbols: [/[pP]/, /[oO]/, /[iI]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'spatial_data_type$subexpression$1',
          symbols: ['spatial_data_type$subexpression$1$subexpression$2']
        }, {
          name: 'spatial_data_type$subexpression$1$subexpression$3',
          symbols: [/[lL]/, /[iI]/, /[nN]/, /[eE]/, /[sS]/, /[tT]/, /[rR]/, /[iI]/, /[nN]/, /[gG]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'spatial_data_type$subexpression$1',
          symbols: ['spatial_data_type$subexpression$1$subexpression$3']
        }, {
          name: 'spatial_data_type$subexpression$1$subexpression$4',
          symbols: [/[pP]/, /[oO]/, /[lL]/, /[yY]/, /[gG]/, /[oO]/, /[nN]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'spatial_data_type$subexpression$1',
          symbols: ['spatial_data_type$subexpression$1$subexpression$4']
        }, {
          name: 'spatial_data_type$subexpression$1$subexpression$5',
          symbols: [/[mM]/, /[uU]/, /[lL]/, /[tT]/, /[iI]/, /[pP]/, /[oO]/, /[iI]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'spatial_data_type$subexpression$1',
          symbols: ['spatial_data_type$subexpression$1$subexpression$5']
        }, {
          name: 'spatial_data_type$subexpression$1$subexpression$6',
          symbols: [/[mM]/, /[uU]/, /[lL]/, /[tT]/, /[iI]/, /[lL]/, /[iI]/, /[nN]/, /[eE]/, /[sS]/, /[tT]/, /[rR]/, /[iI]/, /[nN]/, /[gG]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'spatial_data_type$subexpression$1',
          symbols: ['spatial_data_type$subexpression$1$subexpression$6']
        }, {
          name: 'spatial_data_type$subexpression$1$subexpression$7',
          symbols: [/[mM]/, /[uU]/, /[lL]/, /[tT]/, /[iI]/, /[pP]/, /[oO]/, /[lL]/, /[yY]/, /[gG]/, /[oO]/, /[nN]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'spatial_data_type$subexpression$1',
          symbols: ['spatial_data_type$subexpression$1$subexpression$7']
        }, {
          name: 'spatial_data_type$subexpression$1$subexpression$8',
          symbols: [/[gG]/, /[eE]/, /[oO]/, /[mM]/, /[eE]/, /[tT]/, /[rR]/, /[yY]/, /[cC]/, /[oO]/, /[lL]/, /[lL]/, /[eE]/, /[cC]/, /[tT]/, /[iI]/, /[oO]/, /[nN]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'spatial_data_type$subexpression$1',
          symbols: ['spatial_data_type$subexpression$1$subexpression$8']
        }, {
          name: 'spatial_data_type',
          symbols: ['spatial_data_type$subexpression$1'],
          postprocess: e => ({ type: e[0][0].toUpperCase() })
        }, {
          name: 'column_def$ebnf$1$subexpression$1',
          symbols: ['__', 'column_def_options']
        }, {
          name: 'column_def$ebnf$1',
          symbols: ['column_def$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'column_def$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'column_def',
          symbols: ['_', 'field_name', '__', 'data_type', 'column_def$ebnf$1', '_'],
          postprocess: e => {
            let s = !0; let n = !1; let t = ''; let o = null; let r = !1; let i = ''; let p = ''
            if (e[4] && e[4][1]) {
              for (const $ of e[4][1]) {
                const e = $
                if (e && e.key) {
                  switch (e.key) {
                    case 'NOT NULL':
                      s = !e.value
                      break
                    case 'AUTO_INCREMENT':
                      n = e.value
                      break
                    case 'COMMENT':
                      t = e.value
                      break
                    case 'DEFAULT':
                      o = e.value
                      break
                    case 'ON UPDATE CURRENT_TIMESTAMP':
                      r = e.value
                      break
                    case 'CHARSET':
                      i = e.value
                      break
                    case 'COLLATE':
                      p = e.value
                  }
                }
              }
            }
            const $ = { type: 'column', name: e[1], data_type: e[3], allow_null: s, comment: t }
            return n && ($.auto_increment = n), o !== null && ($.default_value = o), r && ($.on_update_current_timestamp = r), i && ($.charset = i), p && ($.collate = p), $
          }
        }, {
          name: 'column_def_options',
          symbols: ['_', 'column_def_option', '_'],
          postprocess: e => [e[1]]
        }, {
          name: 'column_def_options',
          symbols: ['column_def_options', '__', 'column_def_option'],
          postprocess: e => e[0].concat(e[2])
        }, {
          name: 'column_def_option',
          symbols: ['field_not_null'],
          postprocess: e => e[0]
        }, {
          name: 'column_def_option',
          symbols: ['field_auto_increment'],
          postprocess: e => e[0]
        }, {
          name: 'column_def_option',
          symbols: ['field_comment'],
          postprocess: e => e[0]
        }, {
          name: 'column_def_option',
          symbols: ['field_default_value'],
          postprocess: e => e[0]
        }, {
          name: 'column_def_option',
          symbols: ['field_update_value'],
          postprocess: e => e[0]
        }, {
          name: 'column_def_option',
          symbols: ['field_charset'],
          postprocess: e => e[0]
        }, {
          name: 'field_not_null$subexpression$1$subexpression$1',
          symbols: [/[nN]/, /[oO]/, /[tT]/, { literal: ' ' }, /[nN]/, /[uU]/, /[lL]/, /[lL]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_not_null$subexpression$1',
          symbols: ['field_not_null$subexpression$1$subexpression$1']
        }, {
          name: 'field_not_null$subexpression$1$subexpression$2',
          symbols: [/[nN]/, /[uU]/, /[lL]/, /[lL]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_not_null$subexpression$1',
          symbols: ['field_not_null$subexpression$1$subexpression$2']
        }, {
          name: 'field_not_null',
          symbols: ['field_not_null$subexpression$1'],
          postprocess: e => ({ key: 'NOT NULL', value: e[0][0] && e[0][0].toUpperCase() === 'NOT NULL' })
        }, {
          name: 'field_auto_increment$subexpression$1',
          symbols: [/[aA]/, /[uU]/, /[tT]/, /[oO]/, { literal: '_' }, /[iI]/, /[nN]/, /[cC]/, /[rR]/, /[eE]/, /[mM]/, /[eE]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_auto_increment',
          symbols: ['field_auto_increment$subexpression$1'],
          postprocess: e => ({ key: 'AUTO_INCREMENT', value: !0 })
        }, {
          name: 'field_default_value$subexpression$1',
          symbols: [/[dD]/, /[eE]/, /[fF]/, /[aA]/, /[uU]/, /[lL]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_default_value$subexpression$2',
          symbols: ['decimal']
        }, {
          name: 'field_default_value$subexpression$2',
          symbols: ['string']
        }, {
          name: 'field_default_value$subexpression$2$subexpression$1',
          symbols: [/[cC]/, /[uU]/, /[rR]/, /[rR]/, /[eE]/, /[nN]/, /[tT]/, { literal: '_' }, /[tT]/, /[iI]/, /[mM]/, /[eE]/, /[sS]/, /[tT]/, /[aA]/, /[mM]/, /[pP]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_default_value$subexpression$2',
          symbols: ['field_default_value$subexpression$2$subexpression$1']
        }, {
          name: 'field_default_value$subexpression$2$subexpression$2',
          symbols: [/[nN]/, /[uU]/, /[lL]/, /[lL]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_default_value$subexpression$2',
          symbols: ['field_default_value$subexpression$2$subexpression$2']
        }, {
          name: 'field_default_value',
          symbols: ['field_default_value$subexpression$1', '__', 'field_default_value$subexpression$2'],
          postprocess: e => ({ key: 'DEFAULT', value: e[2][0] })
        }, {
          name: 'field_update_value$subexpression$1',
          symbols: [/[oO]/, /[nN]/, { literal: ' ' }, /[uU]/, /[pP]/, /[dD]/, /[aA]/, /[tT]/, /[eE]/, { literal: ' ' }, /[cC]/, /[uU]/, /[rR]/, /[rR]/, /[eE]/, /[nN]/, /[tT]/, { literal: '_' }, /[tT]/, /[iI]/, /[mM]/, /[eE]/, /[sS]/, /[tT]/, /[aA]/, /[mM]/, /[pP]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_update_value',
          symbols: ['field_update_value$subexpression$1'],
          postprocess: e => ({ key: 'ON UPDATE CURRENT_TIMESTAMP', value: !0 })
        }, {
          name: 'field_charset$subexpression$1$subexpression$1',
          symbols: [/[cC]/, /[hH]/, /[aA]/, /[rR]/, /[aA]/, /[cC]/, /[tT]/, /[eE]/, /[rR]/, { literal: ' ' }, /[sS]/, /[eE]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_charset$subexpression$1',
          symbols: ['field_charset$subexpression$1$subexpression$1']
        }, {
          name: 'field_charset$subexpression$1$subexpression$2',
          symbols: [/[cC]/, /[hH]/, /[aA]/, /[rR]/, /[sS]/, /[eE]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_charset$subexpression$1',
          symbols: ['field_charset$subexpression$1$subexpression$2']
        }, {
          name: 'field_charset',
          symbols: ['field_charset$subexpression$1', '__', 'word'],
          postprocess: e => ({ key: 'CHARSET', value: e[2] })
        }, {
          name: 'field_charset$subexpression$2',
          symbols: [/[cC]/, /[oO]/, /[lL]/, /[lL]/, /[aA]/, /[tT]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_charset',
          symbols: ['field_charset$subexpression$2', '__', 'word'],
          postprocess: e => ({ key: 'COLLATE', value: e[2] })
        }, {
          name: 'field_comment$subexpression$1',
          symbols: [/[cC]/, /[oO]/, /[mM]/, /[mM]/, /[eE]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'field_comment',
          symbols: ['field_comment$subexpression$1', '__', 'string'],
          postprocess: e => ({ key: 'COMMENT', value: e[2] })
        }, {
          name: 'create_definition_list',
          symbols: ['create_definition'],
          postprocess: e => [e[0]]
        }, {
          name: 'create_definition_list',
          symbols: ['create_definition_list', { literal: ',' }, 'create_definition'],
          postprocess: e => e[0].concat(e[2])
        }, { name: 'create_definition', symbols: ['column_def'], postprocess: e => e[0] }, {
          name: 'create_definition',
          symbols: ['primary_key'],
          postprocess: e => e[0]
        }, { name: 'create_definition', symbols: ['index_key'], postprocess: e => e[0] }, {
          name: 'create_definition',
          symbols: ['unique_key'],
          postprocess: e => e[0]
        }, {
          name: 'table_options',
          symbols: ['_', 'table_option', '_'],
          postprocess: e => [e[1]]
        }, {
          name: 'table_options',
          symbols: ['table_options', '__', 'table_option'],
          postprocess: e => e[0].concat(e[2])
        }, {
          name: 'table_option',
          symbols: ['table_option_engine'],
          postprocess: e => ({ key: 'ENGINE', value: e[0] })
        }, {
          name: 'table_option',
          symbols: ['table_option_auto_increment'],
          postprocess: e => ({ key: 'AUTO_INCREMENT', value: e[0] })
        }, {
          name: 'table_option',
          symbols: ['table_option_charset'],
          postprocess: e => ({ key: 'CHARSET', value: e[0] })
        }, {
          name: 'table_option',
          symbols: ['table_option_collate'],
          postprocess: e => ({ key: 'COLLATE', value: e[0] })
        }, {
          name: 'table_option',
          symbols: ['table_option_comment'],
          postprocess: e => ({ key: 'COMMENT', value: e[0] })
        }, {
          name: 'table_option_engine$subexpression$1',
          symbols: [/[eE]/, /[nN]/, /[gG]/, /[iI]/, /[nN]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'table_option_engine',
          symbols: ['table_option_engine$subexpression$1', '_', { literal: '=' }, '_', 'word'],
          postprocess: e => e[4]
        }, {
          name: 'table_option_auto_increment$subexpression$1',
          symbols: [/[aA]/, /[uU]/, /[tT]/, /[oO]/, { literal: '_' }, /[iI]/, /[nN]/, /[cC]/, /[rR]/, /[eE]/, /[mM]/, /[eE]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'table_option_auto_increment',
          symbols: ['table_option_auto_increment$subexpression$1', '_', { literal: '=' }, '_', 'int'],
          postprocess: e => e[4]
        }, {
          name: 'table_option_charset$ebnf$1$subexpression$1$subexpression$1',
          symbols: [/[dD]/, /[eE]/, /[fF]/, /[aA]/, /[uU]/, /[lL]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'table_option_charset$ebnf$1$subexpression$1',
          symbols: ['table_option_charset$ebnf$1$subexpression$1$subexpression$1', '__']
        }, {
          name: 'table_option_charset$ebnf$1',
          symbols: ['table_option_charset$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'table_option_charset$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'table_option_charset$subexpression$1$subexpression$1',
          symbols: [/[cC]/, /[hH]/, /[aA]/, /[rR]/, /[sS]/, /[eE]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'table_option_charset$subexpression$1',
          symbols: ['table_option_charset$subexpression$1$subexpression$1']
        }, {
          name: 'table_option_charset$subexpression$1$subexpression$2',
          symbols: [/[cC]/, /[hH]/, /[aA]/, /[rR]/, /[aA]/, /[cC]/, /[tT]/, /[eE]/, /[rR]/, { literal: ' ' }, /[sS]/, /[eE]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'table_option_charset$subexpression$1',
          symbols: ['table_option_charset$subexpression$1$subexpression$2']
        }, {
          name: 'table_option_charset',
          symbols: ['table_option_charset$ebnf$1', 'table_option_charset$subexpression$1', '_', { literal: '=' }, '_', 'word'],
          postprocess: e => e[5]
        }, {
          name: 'table_option_collate$ebnf$1$subexpression$1$subexpression$1',
          symbols: [/[dD]/, /[eE]/, /[fF]/, /[aA]/, /[uU]/, /[lL]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'table_option_collate$ebnf$1$subexpression$1',
          symbols: ['table_option_collate$ebnf$1$subexpression$1$subexpression$1', '__']
        }, {
          name: 'table_option_collate$ebnf$1',
          symbols: ['table_option_collate$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'table_option_collate$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'table_option_collate$subexpression$1',
          symbols: [/[cC]/, /[oO]/, /[lL]/, /[lL]/, /[aA]/, /[tT]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'table_option_collate',
          symbols: ['table_option_collate$ebnf$1', 'table_option_collate$subexpression$1', '_', { literal: '=' }, '_', 'word'],
          postprocess: e => e[5]
        }, {
          name: 'table_option_comment$subexpression$1',
          symbols: [/[cC]/, /[oO]/, /[mM]/, /[mM]/, /[eE]/, /[nN]/, /[tT]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'table_option_comment',
          symbols: ['table_option_comment$subexpression$1', '_', { literal: '=' }, '_', 'string'],
          postprocess: e => e[4]
        }, {
          name: 'index_type$subexpression$1',
          symbols: [/[uU]/, /[sS]/, /[iI]/, /[nN]/, /[gG]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'index_type$subexpression$2$subexpression$1',
          symbols: [/[bB]/, /[tT]/, /[rR]/, /[eE]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'index_type$subexpression$2',
          symbols: ['index_type$subexpression$2$subexpression$1']
        }, {
          name: 'index_type$subexpression$2$subexpression$2',
          symbols: [/[hH]/, /[aA]/, /[sS]/, /[hH]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'index_type$subexpression$2',
          symbols: ['index_type$subexpression$2$subexpression$2']
        }, {
          name: 'index_type',
          symbols: ['_', 'index_type$subexpression$1', '_', 'index_type$subexpression$2', '_'],
          postprocess: e => e[3]
        }, {
          name: 'primary_key$subexpression$1',
          symbols: [/[pP]/, /[rR]/, /[iI]/, /[mM]/, /[aA]/, /[rR]/, /[yY]/, { literal: ' ' }, /[kK]/, /[eE]/, /[yY]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, { name: 'primary_key$ebnf$1', symbols: ['index_type'], postprocess: s }, {
          name: 'primary_key$ebnf$1',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, { name: 'primary_key$ebnf$2', symbols: [{ literal: ',' }], postprocess: s }, {
          name: 'primary_key$ebnf$2',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, {
          name: 'primary_key',
          symbols: ['_', 'primary_key$subexpression$1', '_', { literal: '(' }, 'key_field_list', { literal: ')' }, 'primary_key$ebnf$1', '_', 'primary_key$ebnf$2'],
          postprocess: e => ({ type: 'primary_key', fields: e[4] })
        }, {
          name: 'index_key$subexpression$1', symbols: [/[kK]/, /[eE]/, /[yY]/], postprocess: function(e) {
            return e.join('')
          }
        }, { name: 'index_key$ebnf$1', symbols: ['index_type'], postprocess: s }, {
          name: 'index_key$ebnf$1',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, { name: 'index_key$ebnf$2', symbols: [{ literal: ',' }], postprocess: s }, {
          name: 'index_key$ebnf$2',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, {
          name: 'index_key',
          symbols: ['_', 'index_key$subexpression$1', '__', 'field_name', '_', { literal: '(' }, 'key_field_list', { literal: ')' }, 'index_key$ebnf$1', '_', 'index_key$ebnf$2'],
          postprocess: e => ({ type: 'index_key', name: e[3], fields: e[6] })
        }, {
          name: 'unique_key$subexpression$1',
          symbols: [/[uU]/, /[nN]/, /[iI]/, /[qQ]/, /[uU]/, /[eE]/, { literal: ' ' }, /[kK]/, /[eE]/, /[yY]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, { name: 'unique_key$ebnf$1', symbols: ['index_type'], postprocess: s }, {
          name: 'unique_key$ebnf$1',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, { name: 'unique_key$ebnf$2', symbols: [{ literal: ',' }], postprocess: s }, {
          name: 'unique_key$ebnf$2',
          symbols: [],
          postprocess: function(e) {
            return null
          }
        }, {
          name: 'unique_key',
          symbols: ['_', 'unique_key$subexpression$1', '__', 'field_name', '_', { literal: '(' }, 'key_field_list', { literal: ')' }, 'unique_key$ebnf$1', '_', 'unique_key$ebnf$2'],
          postprocess: e => ({ type: 'unique_key', name: e[3], fields: e[6] })
        }, {
          name: 'key_field_list$ebnf$1$subexpression$1',
          symbols: [{ literal: '(' }, '_', 'unsigned_int', '_', { literal: ')' }]
        }, {
          name: 'key_field_list$ebnf$1',
          symbols: ['key_field_list$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'key_field_list$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'key_field_list',
          symbols: ['_', 'field_name', 'key_field_list$ebnf$1', '_'],
          postprocess: e => [e[1]]
        }, {
          name: 'key_field_list',
          symbols: ['key_field_list', '_', { literal: ',' }, '_', 'field_name'],
          postprocess: e => e[0].concat(e[4])
        }, { name: 'comment', symbols: ['inline_comment'], postprocess: e => e[0] }, {
          name: 'comment',
          symbols: ['multiline_comment'],
          postprocess: e => e[0]
        }, { name: 'comment', symbols: ['newline'], postprocess: e => null }, {
          name: 'inline_comment$subexpression$1',
          symbols: [{ literal: '#' }]
        }, {
          name: 'inline_comment$subexpression$1$string$1',
          symbols: [{ literal: '-' }, { literal: '-' }],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'inline_comment$subexpression$1',
          symbols: ['inline_comment$subexpression$1$string$1']
        }, { name: 'inline_comment$ebnf$1', symbols: [] }, {
          name: 'inline_comment$ebnf$1',
          symbols: ['inline_comment$ebnf$1', /[^\n]/],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'inline_comment',
          symbols: ['space', 'inline_comment$subexpression$1', 'inline_comment$ebnf$1', /[\n]/],
          postprocess: e => e[2].join('')
        }, {
          name: 'multiline_comment$string$1',
          symbols: [{ literal: '/' }, { literal: '*' }],
          postprocess: function(e) {
            return e.join('')
          }
        }, { name: 'multiline_comment$ebnf$1', symbols: [] }, {
          name: 'multiline_comment$ebnf$1$subexpression$1',
          symbols: [/[^\*]/]
        }, {
          name: 'multiline_comment$ebnf$1$subexpression$1',
          symbols: [/[\*]/, /[^\/\*]/]
        }, {
          name: 'multiline_comment$ebnf$1',
          symbols: ['multiline_comment$ebnf$1', 'multiline_comment$ebnf$1$subexpression$1'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'multiline_comment$string$2',
          symbols: [{ literal: '*' }, { literal: '/' }],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'multiline_comment',
          symbols: ['space', 'multiline_comment$string$1', 'multiline_comment$ebnf$1', 'multiline_comment$string$2', 'space'],
          postprocess: e => [].concat(...e[2]).join('')
        }, { name: 'main', symbols: ['blocks'], postprocess: e => e[0] }, {
          name: 'blocks$ebnf$1',
          symbols: []
        }, {
          name: 'blocks$ebnf$1', symbols: ['blocks$ebnf$1', { literal: ';' }], postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, { name: 'blocks', symbols: ['block', 'blocks$ebnf$1'], postprocess: e => e[0] }, {
          name: 'blocks$ebnf$2',
          symbols: [{ literal: ';' }]
        }, {
          name: 'blocks$ebnf$2', symbols: ['blocks$ebnf$2', { literal: ';' }], postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, { name: 'blocks$ebnf$3', symbols: [] }, {
          name: 'blocks$ebnf$3',
          symbols: ['blocks$ebnf$3', { literal: ';' }],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'blocks',
          symbols: ['blocks', 'blocks$ebnf$2', 'block', 'blocks$ebnf$3'],
          postprocess: e => (e[0] || []).concat(e[2] || [])
        }, { name: 'block', symbols: ['space'], postprocess: e => null }, {
          name: 'block$ebnf$1',
          symbols: ['comment']
        }, {
          name: 'block$ebnf$1', symbols: ['block$ebnf$1', 'comment'], postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, { name: 'block', symbols: ['block$ebnf$1'], postprocess: e => null }, {
          name: 'block$ebnf$2',
          symbols: ['statement']
        }, {
          name: 'block$ebnf$2', symbols: ['block$ebnf$2', 'statement'], postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'block',
          symbols: ['block$ebnf$2'],
          postprocess: e => e[0].filter(e => e)
        }, { name: 'statement$ebnf$1', symbols: [] }, {
          name: 'statement$ebnf$1',
          symbols: ['statement$ebnf$1', 'comment'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, { name: 'statement$ebnf$2', symbols: [] }, {
          name: 'statement$ebnf$2',
          symbols: ['statement$ebnf$2', 'comment'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'statement',
          symbols: ['statement$ebnf$1', 'P_CREATE_TABLE', 'statement$ebnf$2'],
          postprocess: e => e[1]
        }, { name: 'statement$ebnf$3', symbols: [] }, {
          name: 'statement$ebnf$3',
          symbols: ['statement$ebnf$3', 'comment'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, { name: 'statement$ebnf$4', symbols: [] }, {
          name: 'statement$ebnf$4',
          symbols: ['statement$ebnf$4', 'comment'],
          postprocess: function(e) {
            return e[0].concat([e[1]])
          }
        }, {
          name: 'statement',
          symbols: ['statement$ebnf$3', 'P_SET', 'statement$ebnf$4'],
          postprocess: e => null
        }, {
          name: 'P_CREATE_TABLE$subexpression$1',
          symbols: [/[cC]/, /[rR]/, /[eE]/, /[aA]/, /[tT]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'P_CREATE_TABLE$ebnf$1$subexpression$1$subexpression$1',
          symbols: [/[tT]/, /[eE]/, /[mM]/, /[pP]/, /[oO]/, /[rR]/, /[aA]/, /[rR]/, /[yY]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'P_CREATE_TABLE$ebnf$1$subexpression$1',
          symbols: ['__', 'P_CREATE_TABLE$ebnf$1$subexpression$1$subexpression$1']
        }, {
          name: 'P_CREATE_TABLE$ebnf$1',
          symbols: ['P_CREATE_TABLE$ebnf$1$subexpression$1'],
          postprocess: s
        }, {
          name: 'P_CREATE_TABLE$ebnf$1', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'P_CREATE_TABLE$subexpression$2',
          symbols: [/[tT]/, /[aA]/, /[bB]/, /[lL]/, /[eE]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'P_CREATE_TABLE$ebnf$2$subexpression$1$subexpression$1',
          symbols: [/[iI]/, /[fF]/, { literal: ' ' }, /[nN]/, /[oO]/, /[tT]/, { literal: ' ' }, /[eE]/, /[xX]/, /[iI]/, /[sS]/, /[tT]/, /[sS]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'P_CREATE_TABLE$ebnf$2$subexpression$1',
          symbols: ['__', 'P_CREATE_TABLE$ebnf$2$subexpression$1$subexpression$1']
        }, {
          name: 'P_CREATE_TABLE$ebnf$2',
          symbols: ['P_CREATE_TABLE$ebnf$2$subexpression$1'],
          postprocess: s
        }, {
          name: 'P_CREATE_TABLE$ebnf$2', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'P_CREATE_TABLE$ebnf$3',
          symbols: ['table_options'],
          postprocess: s
        }, {
          name: 'P_CREATE_TABLE$ebnf$3', symbols: [], postprocess: function(e) {
            return null
          }
        }, {
          name: 'P_CREATE_TABLE',
          symbols: ['_', 'P_CREATE_TABLE$subexpression$1', 'P_CREATE_TABLE$ebnf$1', '__', 'P_CREATE_TABLE$subexpression$2', 'P_CREATE_TABLE$ebnf$2', '__', 'table_name', '_', { literal: '(' }, 'create_definition_list', { literal: ')' }, 'P_CREATE_TABLE$ebnf$3', '_'],
          postprocess: e => ({ type: 'create_table', name: e[7], columns: e[10], options: e[12] })
        }, {
          name: 'P_SET$subexpression$1', symbols: [/[sS]/, /[eE]/, /[tT]/], postprocess: function(e) {
            return e.join('')
          }
        }, {
          name: 'P_SET$subexpression$2$subexpression$1',
          symbols: [/[nN]/, /[aA]/, /[mM]/, /[eE]/, /[sS]/],
          postprocess: function(e) {
            return e.join('')
          }
        }, { name: 'P_SET$subexpression$2', symbols: ['P_SET$subexpression$2$subexpression$1'] }, {
          name: 'P_SET',
          symbols: ['_', 'P_SET$subexpression$1', '__', 'P_SET$subexpression$2', '__', 'set_value', '_'],
          postprocess: e => null
        }],
        ParserStart: 'main'
      }
      e.exports = n
    }())
  })
  var r = {
    parse(e) {
      const s = new t.Parser(t.Grammar.fromCompiled(o))
      return s.feed(e), s.results[0]
    }
  }; var i = r.parse
  e.default = r, e.parse = i, Object.defineProperty(e, '__esModule', { value: !0 })
  SqlParser = e
}))

module.exports = {
  SqlParser
}
