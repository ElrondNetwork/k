// Copyright (c) 2013-2019 K Team. All Rights Reserved.
module RULES
  rule
    <k>
    return (x * x + y * y => V:Float);
    ...</k>
    <env>...
      x |-> LX:Int
      y |-> LY:Int
    ...</env>
    <store>...
      LX |-> X:Float
      LY |-> Y:Float
    ...</store>
  requires -100.0 <=Float X <=Float 100.0 andBool -100.0 <=Float Y <=Float 100.0
  ensures -0x1.p-37 <=Float V -Float (X *Float X +Float Y *Float Y) <=Float 0x1.p-37
endmodule

