# WishTournament Bracket System

This document provides an in-depth view of the bracket system used for handling matchups, rounds, and games within the **WishTournament** app.

---

## Bracket Overview

A **bracket** represents the structure of matchups between teams or players in a tournament. For this system, we'll focus on the **single elimination** bracket as the default, but it can be extended to support other formats (e.g., double elimination, round-robin).

### Key Elements:
1. **Bracket**: The overall structure of the tournament.
2. **Rounds**: Stages within the bracket where multiple games take place.
3. **Games**: Individual matchups between two teams or players within a round.
4. **Progression**: Winners from each game move forward to the next round.

---

## Bracket Data Models

### Bracket Model

```go
type Bracket struct {
    ID           uuid.UUID `json:"id"`               // UUID for the bracket
    TournamentID uuid.UUID `json:"tournament_id"`    // UUID of the associated tournament
    Rounds       []Round   `json:"rounds"`           // array of rounds in the bracket
    Type         string    `json:"type"`             // e.g., "single_elimination"
    CreatedAt    time.Time `json:"created_at"`
}
```

- **TournamentID**: Links the bracket to a specific tournament.
- **Rounds**: An array of **Round** objects, representing each stage of the bracket.
- **Type**: The type of bracket (e.g., single elimination).

### Round Model

Each bracket is divided into multiple rounds, with each round consisting of multiple games.

```go
type Round struct {
    ID         uuid.UUID  `json:"id"`               // UUID for the round
    BracketID  uuid.UUID  `json:"bracket_id"`       // UUID of the bracket this round belongs to
    RoundNumber int       `json:"round_number"`     // 1 = quarter-finals, 2 = semi-finals, etc.
    Games      []Game     `json:"games"`            // list of games in this round
    CreatedAt  time.Time  `json:"created_at"`
}
```

- **RoundNumber**: Represents the stage of the tournament (e.g., Round 1, Quarter-finals).
- **Games**: An array of **Game** objects, representing each game played in this round.

### Game Model (Bracket-Specific)

A game represents a match between two teams in a specific round.

```go
type Game struct {
    ID         uuid.UUID  `json:"id"`                // UUID for the game
    RoundID    uuid.UUID  `json:"round_id"`          // UUID for the round this game belongs to
    Team1ID    uuid.UUID  `json:"team1_id"`          // UUID for the first team
    Team2ID    uuid.UUID  `json:"team2_id"`          // UUID for the second team
    Team1Score float64    `json:"team1_score"`       // score for team 1
    Team2Score float64    `json:"team2_score"`       // score for team 2
    WinnerID   *uuid.UUID `json:"winner_id"`         // UUID of the winning team (nullable if undecided)
    IsTie      bool       `json:"is_tie"`            // true if the game resulted in a tie
    CreatedAt  time.Time  `json:"created_at"`
}
```

- **RoundID**: Links the game to a specific round.
- **WinnerID**: Stores the **winning team's** UUID. This is used to progress the winner to the next round.
- **IsTie**: Indicates if the game ended in a tie (in cases where ties are allowed).

---

## Bracket API Endpoints

### 1. **Create a Bracket**

Creates a new bracket for a specific tournament. The bracket will be initialized with Round 1 based on the participating teams.

```http
POST /tournaments/:id/brackets
```

- **Request Body**:
   ```json
   {
     "type": "single_elimination"
   }
   ```

- **Response**:
   ```json
   {
     "id": "123e4567-e89b-12d3-a456-426614174000",
     "tournament_id": "123e4567-e89b-12d3-a456-426614174111",
     "type": "single_elimination",
     "rounds": []
   }
   ```

### 2. **Get a Bracket for a Tournament**

Retrieves the current state of a bracket, including all rounds and games.

```http
GET /tournaments/:id/brackets
```

- **Response**:
   ```json
   {
     "id": "123e4567-e89b-12d3-a456-426614174000",
     "tournament_id": "123e4567-e89b-12d3-a456-426614174111",
     "type": "single_elimination",
     "rounds": [
       {
         "id": "round-uuid-1",
         "round_number": 1,
         "games": [
           {
             "id": "game-uuid-1",
             "team1_id": "team1-uuid",
             "team2_id": "team2-uuid",
             "team1_score": 100,
             "team2_score": 98,
             "winner_id": "team1-uuid",
             "is_tie": false
           }
         ]
       }
     ]
   }
   ```

### 3. **Update Game Results (Set Winner)**

After a game is played, update the scores and the winner. This endpoint allows you to set the winner for a specific game in a round.

```http
PUT /brackets/:bracketId/games/:gameId
```

- **Request Body**:
   ```json
   {
     "team1_score": 80,
     "team2_score": 85,
     "winner_id": "team2-uuid",
     "is_tie": false
   }
   ```

- **Response**:
   ```json
   {
     "id": "game-uuid",
     "team1_score": 80,
     "team2_score": 85,
     "winner_id": "team2-uuid",
     "is_tie": false
   }
   ```

### 4. **Move to Next Round**

Progress the bracket to the next round once all games in the current round are completed.

```http
POST /brackets/:bracketId/rounds/:roundId/next
```

- **Response**:
   ```json
   {
     "round_number": 2,
     "games": [
       {
         "id": "game-uuid-new",
         "team1_id": "winner-from-previous-round",
         "team2_id": "another-winner"
       }
     ]
   }
   ```

---

## Single Elimination Flow

Here’s an overview of how the single elimination bracket works:

1. **Create Bracket**: When a tournament is started, the system creates the bracket and populates Round 1 with matchups based on the participating teams.
    - Teams are randomly paired, or seeded, based on the tournament's configuration.

2. **Record Results**: As games are played, the host records the **scores** and **winner** using the `PUT /brackets/:bracketId/games/:gameId` endpoint.

3. **Progress to Next Round**: Once all games in the current round are completed, the winners progress to the next round via the `POST /brackets/:bracketId/rounds/:roundId/next` endpoint.
    - The system automatically generates the matchups for the next round based on the previous round’s winners.

4. **Final Round**: This process repeats until a final winner is determined. The last round will have a single game, and the winner of that game is the tournament champion.

---

## Advanced Bracket Features (Future)

- **Double Elimination**: A more complex system where losing teams get a second chance to compete in a **losers' bracket**.
- **Round Robin**: A system where each team plays against every other team, with the overall winner determined by the total number of wins.
- **Seeding**: Allow the host to seed teams based on their rankings or past performances, instead of random pairings.

---

## Example Flow (Single Elimination Bracket)

1. **Bracket Creation**:
    - Teams A, B, C, and D enter a single elimination tournament.
    - **Round 1**: Team A plays Team B, Team C plays Team D.
    - Winners progress to Round 2.

2. **Round 1 Results**:
    - Team A beats Team B.
    - Team D beats Team C.

3. **Round 2 (Final)**:
    - Team A plays Team D. The winner of this game is crowned the tournament champion.

---