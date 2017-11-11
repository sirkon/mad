grammar Raw;


set
    : (line eoc)* (line EOF)?
    |
    ;

eoc
    : NEWLINE*;


line
    : (IDENTIFIER '=' value)
    ;

value
    : STRING_LITERAL
    | NUMBER
    | boolean
    | INLINE_STRING
    | IDENTIFIER
    ;

boolean
    : 'true'
    | 'false';

IDENTIFIER
    : [a-zA-Z_] [a-zA-Z0-9_]*
    ;

STRING_LITERAL
    : '"' (~["\\\r\n] | EscapeSequence)* '"'
    ;


NUMBER
    : '-'? INT ('.' [0-9] +)? EXP?
    ;

INT
   : '0' | [1-9] [0-9]*
   ;
// no leading zeros
fragment EXP
   : [Ee] [+\-]? INT
   ;

fragment EscapeSequence
    : '\\' [btnfr"'\\]
    ;

INLINE_STRING
    : ~(' ' | '\r' | '\t' | '\n' | '=')+
    ;


NEWLINE: '\n';

WS
    :
    [ \t\r] -> skip ;


