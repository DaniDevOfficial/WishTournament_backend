# WishTournament Backend

This is the backend for the WishTournament app, which tracks teams and players in tournaments.

---

## Use Cases

### User
1. As a user, I want to CRUD an account.
2. As a user, I want to CRUD a tournament.
3. As a user, I want to get invited to a tournament.
4. As a user, I want to leave a tournament.
5. As a user, I want to follow other users (like Instagram or GitHub).

### Hoster (Creator of a Tournament)
1. As a hoster, I want to create teams with players in them.
2. As a hoster, I want to create brackets for teams.
3. As a hoster, I want to set points and a winner for a single game.
4. As a hoster, I want to CRUD notes for a specific game.
5. As a hoster, I want to add either real or anonymous players to a team and track them correctly.
6. As a hoster, I want to define the type of points for a tournament (e.g., points in basketball, time in running, length in jumping, etc.) and whether the highest or lowest score wins.
7. As a hoster, I want to allow ties in a game.
8. As a hoster, I want to have a leaderboard for the entire tournament.
9. As a hoster, I want to have a **pregame phase** where tournament details can be changed until the tournament starts, after which no more changes can be made.

---

## Backend API Endpoints

### User-Related Endpoints

| **Use Case**                | **Endpoint**                     | **Method** |
|-----------------------------|----------------------------------|------------|
| CRUD an Account             | `/users`                         | POST, GET, PUT, DELETE |
| CRUD a Tournament           | `/tournaments`                   | POST, GET, PUT, DELETE |
| Get invited to a Tournament | `/tournaments/:id/invite`        | POST |
| Leave a Tournament          | `/tournaments/:id/leave`         | POST |
| Follow other users          | `/users/:id/follow`              | POST |

### Hoster-Related Endpoints

| **Use Case**                                           | **Endpoint**                             | **Method** |
|--------------------------------------------------------|------------------------------------------|------------|
| Create teams                                           | `/tournaments/:id/teams`                 | POST |
| Create brackets for teams                              | `/tournaments/:id/brackets`              | POST |
| Set points and winner for a game                       | `/tournaments/:id/games/:gameId`         | PUT |
| CRUD notes for a game                                  | `/tournaments/:id/games/:gameId/notes`   | POST, PUT, DELETE |
| Add real/anonymous players to a team                   | `/tournaments/:id/teams/:teamId/players` | POST |
| Get leaderboard                                        | `/tournaments/:id/leaderboard`           | GET |

---

## Data Models (with UUIDs and Start Date/Time)

### User Model

```go
import "github.com/google/uuid"

type User struct {
    ID        uuid.UUID `json:"id"`              // UUID for global uniqueness
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`
    Following []uuid.UUID `json:"following"`     // list of user UUIDs
    CreatedAt time.Time `json:"created_at"`
}
```

### Tournament Model

```go
type Tournament struct {
    ID            uuid.UUID `json:"id"`              // UUID for the tournament
    Name          string    `json:"name"`
    HostID        uuid.UUID `json:"host_id"`         // UUID of the user who created the tournament
    Teams         []Team    `json:"teams"`           // array of teams participating
    PointType     string    `json:"point_type"`      // "points", "time", "distance"
    ScoringMethod string    `json:"scoring_method"`  // "highest" or "lowest"
    CanTie        bool      `json:"can_tie"`         // allow ties in the tournament
    StartDateTime time.Time `json:"start_date_time"` // when the tournament starts
    IsStarted     bool      `json:"is_started"`      // true if the tournament has started
    CreatedAt     time.Time `json:"created_at"`
}
```

### Team Model

```go
type Team struct {
    ID           uuid.UUID `json:"id"`               // UUID for the team
    Name         string    `json:"name"`
    Players      []Player  `json:"players"`          // players can be real or anonymous
    TournamentID uuid.UUID `json:"tournament_id"`    // UUID of the associated tournament
    CreatedAt    time.Time `json:"created_at"`
}
```

### Player Model

```go
type Player struct {
    ID            uuid.UUID  `json:"id"`                // UUID for the player
    Name          string     `json:"name"`              // Chosen by host for anonymous players
    IsRegistered  bool       `json:"is_registered"`     // True if the player has an account
    UserID        *uuid.UUID `json:"user_id,omitempty"` // Nullable: linked to a User account if registered
    TeamID        uuid.UUID  `json:"team_id"`           // Team the player belongs to
    TournamentID  uuid.UUID  `json:"tournament_id"`     // Tournament they’re participating in
}
```

---

## Scoring Logic

1. When creating a tournament, the host can define:
    - **PointType** (e.g., "points", "time", "distance").
    - **ScoringMethod** ("highest" or "lowest" score wins).
    - **CanTie** (whether ties are allowed).
    - **StartDateTime** (the time when the tournament starts).

2. In a game:
    - If `canTie` is `true` and both teams have equal scores, the game is marked as a **tie** (`isTie: true`).
    - Otherwise, the winner is determined based on the **ScoringMethod** (highest or lowest score).

---

### Pregame Phase & Locking

1. **Pregame Phase**:
    - During the pregame phase (before `startDateTime`), the host can modify tournament details, including teams, players, and configurations.

2. **Locked Phase**:
    - After the tournament’s `startDateTime`, the `isStarted` flag is set to `true`, and no further changes can be made to the tournament settings or teams.
    - Players can no longer be added or removed, and tournament configurations like point types and scoring methods become locked.

---

### Anonymous Players

- Anonymous players can be added without forcing them to create an account. These players will be stored in the **Player** model with:
    - **IsRegistered = false**.
    - No **UserID** (left `null`).
    - A **name** chosen by the tournament host.

Example:

```json
{
  "team": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Team Alpha",
    "players": [
      { "id": "223e4567-e89b-12d3-a456-426614174001", "name": "John Doe", "is_registered": false, "user_id": null },
      { "id": "323e4567-e89b-12d3-a456-426614174002", "name": "RealPlayer1", "is_registered": true, "user_id": "523e4567-e89b-12d3-a456-426614174003" }
    ],
    "tournament_id": "423e4567-e89b-12d3-a456-426614174004"
  }
}
```