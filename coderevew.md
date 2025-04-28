```diff
--- a/src/chess/board.cc
+++ b/src/chess/board.cc
@@ -975,30 +975,42 @@ void ChessBoard::SetFromFen(const std::string& fen, int* no_capture_ply,
     } else {
       throw Exception("Bad fen string: " + fen);
     }
     ++col;
   }

   if (castlings != "-") {
     for (char c : castlings) {
       switch (c) {
         case 'K':
-          castlings_.set_we_can_00();
+          if (our_king_.as_string() == "e1" && our_pieces_.get(0, 7) &&
+              rooks_.get(0, 7)) {
+            castlings_.set_we_can_00();
+          }
           break;
         case 'k':
-          castlings_.set_they_can_00();
+          if (their_king_.as_string() == "e8" && their_pieces_.get(7, 7) &&
+              rooks_.get(7, 7)) {
+            castlings_.set_they_can_00();
+          }
           break;
         case 'Q':
-          castlings_.set_we_can_000();
+          if (our_king_.as_string() == "e1" && our_pieces_.get(0, 0) &&
+              rooks_.get(0, 0)) {
+            castlings_.set_we_can_000();
+          }
           break;
         case 'q':
-          castlings_.set_they_can_000();
+          if (their_king_.as_string() == "e8" && their_pieces_.get(7, 0) &&
+              rooks_.get(7, 0)) {
+            castlings_.set_they_can_000();
+          }
           break;
         default:
           throw Exception("Bad fen string: " + fen);
       }
     }
   }
```

[Description]
The code changes in the `SetFromFen` method relate to setting castling rights based on the FEN string. The original code unconditionally set the castling rights based on the characters 'K', 'k', 'Q', and 'q' in the FEN string. The new code adds checks to ensure that the king and the corresponding rook are in their starting positions before setting the castling rights. Specifically, it checks:

- For 'K' (white kingside castling): The white king must be on "e1", a white piece and a rook must be on "h1" (0,7).
- For 'k' (black kingside castling): The black king must be on "e8", a black piece and a rook must be on "h8" (7,7).
- For 'Q' (white queenside castling): The white king must be on "e1", a white piece and a rook must be on "a1" (0,0).
- For 'q' (black queenside castling): The black king must be on "e8", a black piece and a rook must be on "a8" (7,0).

[Obvious errors]
There are no obvious errors.

[Improvements]

1.  **Readability**: The conditions for castling can be made more readable by extracting them into named boolean variables.

```c++
         case 'K': {
+          bool king_on_e1 = our_king_.as_string() == "e1";
+          bool rook_on_h1 = our_pieces_.get(0, 7) && rooks_.get(0, 7);
+          if (king_on_e1 && rook_on_h1) {
+            castlings_.set_we_can_00();
+          }
+          break;
+        }
+        case 'k': {
+          bool king_on_e8 = their_king_.as_string() == "e8";
+          bool rook_on_h8 = their_pieces_.get(7, 7) && rooks_.get(7, 7);
+          if (king_on_e8 && rook_on_h8) {
+            castlings_.set_they_can_00();
+          }
+          break;
+        }
+        case 'Q': {
+          bool king_on_e1 = our_king_.as_string() == "e1";
+          bool rook_on_a1 = our_pieces_.get(0, 0) && rooks_.get(0, 0);
+          if (king_on_e1 && rook_on_a1) {
+            castlings_.set_we_can_000();
+          }
+          break;
+        }
+        case 'q': {
+          bool king_on_e8 = their_king_.as_string() == "e8";
+          bool rook_on_a8 = their_pieces_.get(7, 0) && rooks_.get(7, 0);
+          if (king_on_e8 && rook_on_a8) {
+            castlings_.set_they_can_000();
+          }
+          break;
+        }
```

2.  **Duplication**: The code repeats the pattern of checking king position and rook presence for each castling right. This duplication can be reduced by creating a helper function.

```c++
 bool CanCastle(bool is_white, bool kingside, const BoardSquare& king_pos,
                 const Bitboard& pieces, const Bitboard& rooks) {
  std::string expected_king_pos = is_white ? "e1" : "e8";
  int rook_row = is_white ? 0 : 7;
  int rook_col = kingside ? 7 : 0;


  return king_pos.as_string() == expected_king_pos && pieces.get(rook_row, rook_col) &&
         rooks.get(rook_row, rook_col);
 }


 // Then, in the switch statement:
 case 'K':
  if (CanCastle(true, true, our_king_, our_pieces_, rooks_)) {
  castlings_.set_we_can_00();
  }
  break;
 case 'k':
  if (CanCastle(false, true, their_king_, their_pieces_, rooks_)) {
  castlings_.set_they_can_00();
  }
  break;
 case 'Q':
  if (CanCastle(true, false, our_king_, our_pieces_, rooks_)) {
  castlings_.set_we_can_000();
  }
  break;
 case 'q':
  if (CanCastle(false, false, their_king_, their_pieces_, rooks_)) {
  castlings_.set_they_can_000();
  }
  break;
```

[Friendly advice]
The added checks are a good step to ensure that castling rights are set correctly based on the board state. The suggestions above are aimed at improving code readability and reducing code duplication. Keep up the good work!
