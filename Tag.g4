grammar Tag;

set
    : tag*;

tag
    : IDENTIFIER ':' STRING_LITERAL ' '*
    ;

IDENTIFIER
    : [0-9a-zA-Z-]+;





STRING_LITERAL
    : '"' (~["])* '"'
    ;



