package services

import (
    "database/sql"
    "fmt"
    "math"
    "rideshare-backend/internal/models"
    "sort"
    "strings"
)

type MatchingService struct {
    db *sql.DB
}

func NewMatchingService(db *sql.DB) *MatchingService {
    return &MatchingService{db: db}
}


// Enhanced FindMatchingTrips with better filtering and scoring
func (s *MatchingService) FindMatchingTrips(userID int, from, to string, maxDistance float64) ([]models.TripMatch, error) {
    // Set minimum similarity threshold to filter out irrelevant results
    minSimilarity := 0.3 // 30% minimum similarity
    
    query := `
        SELECT t.id, t.driver_id, t.from_location, t.to_location, t.departure_time,
               t.max_passengers, t.current_passengers, t.price_per_person, t.description,
               t.status, t.created_at, t.updated_at
        FROM trips t
        WHERE t.status = 'active' 
        AND t.driver_id != $1
        AND t.current_passengers < t.max_passengers
        AND t.departure_time > NOW()
        ORDER BY t.departure_time ASC
    `
    
    rows, err := s.db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query trips: %w", err)
    }
    defer rows.Close()
    
    var allMatches []models.TripMatch
    
    for rows.Next() {
        var trip models.Trip
        
        err := rows.Scan(
            &trip.ID,
            &trip.DriverID,
            &trip.FromLocation,
            &trip.ToLocation,
            &trip.DepartureTime,
            &trip.MaxPassengers,
            &trip.CurrentPassengers,
            &trip.PricePerPerson,
            &trip.Description,
            &trip.Status,
            &trip.CreatedAt,
            &trip.UpdatedAt,
        )
        if err != nil {
            continue
        }
        
        // Calculate comprehensive similarity and distance scores
        similarity := s.calculateEnhancedSimilarity(from, to, trip.FromLocation, trip.ToLocation)
        distance := s.calculateDistance(from, to, trip.FromLocation, trip.ToLocation)
        
        // Include all trips that meet minimum criteria
        if similarity >= minSimilarity && distance <= maxDistance {
            match := models.TripMatch{
                TripID:      trip.ID,
                PassengerID: userID,
                Similarity:  similarity,
                Distance:    distance,
            }
            
            allMatches = append(allMatches, match)
        }
    }
    
    // Sort by similarity (highest first), then by distance (lowest first)
    sort.Slice(allMatches, func(i, j int) bool {
        if math.Abs(allMatches[i].Similarity - allMatches[j].Similarity) < 0.01 {
            // If similarity is very close, prefer lower distance
            return allMatches[i].Distance < allMatches[j].Distance
        }
        return allMatches[i].Similarity > allMatches[j].Similarity
    })
    
    return allMatches, nil
}

// Enhanced similarity calculation with multiple factors
func (s *MatchingService) calculateEnhancedSimilarity(from1, to1, from2, to2 string) float64 {
    // Normalize strings
    from1, to1 = strings.ToLower(strings.TrimSpace(from1)), strings.ToLower(strings.TrimSpace(to1))
    from2, to2 = strings.ToLower(strings.TrimSpace(from2)), strings.ToLower(strings.TrimSpace(to2))
    
    // Calculate individual similarities
    fromSimilarity := s.advancedStringSimilarity(from1, from2)
    toSimilarity := s.advancedStringSimilarity(to1, to2)
    
    // Weight the similarities (destination is slightly more important)
    weightedSimilarity := (fromSimilarity * 0.45) + (toSimilarity * 0.55)
    
    return weightedSimilarity
}

// Advanced string similarity with multiple matching techniques
func (s *MatchingService) advancedStringSimilarity(s1, s2 string) float64 {
    if s1 == s2 {
        return 1.0
    }
    
    if s1 == "" || s2 == "" {
        return 0.0
    }
    
    // Multiple similarity checks
    scores := []float64{
        s.exactMatchScore(s1, s2),
        s.substringMatchScore(s1, s2),
        s.wordMatchScore(s1, s2),
        s.levenshteinSimilarity(s1, s2),
        s.jaroWinklerSimilarity(s1, s2),
    }
    
    // Return the highest score from all methods
    maxScore := 0.0
    for _, score := range scores {
        if score > maxScore {
            maxScore = score
        }
    }
    
    return maxScore
}

func (s *MatchingService) exactMatchScore(s1, s2 string) float64 {
    if s1 == s2 {
        return 1.0
    }
    return 0.0
}

func (s *MatchingService) substringMatchScore(s1, s2 string) float64 {
    if strings.Contains(s1, s2) || strings.Contains(s2, s1) {
        shorter := math.Min(float64(len(s1)), float64(len(s2)))
        longer := math.Max(float64(len(s1)), float64(len(s2)))
        return shorter / longer * 0.9 // 90% for substring matches
    }
    return 0.0
}

func (s *MatchingService) wordMatchScore(s1, s2 string) float64 {
    words1 := strings.Fields(s1)
    words2 := strings.Fields(s2)
    
    if len(words1) == 0 || len(words2) == 0 {
        return 0.0
    }
    
    commonWords := 0
    totalWords := len(words1) + len(words2)
    
    for _, w1 := range words1 {
        for _, w2 := range words2 {
            if w1 == w2 || strings.Contains(w1, w2) || strings.Contains(w2, w1) {
                commonWords++
                break
            }
        }
    }
    
    return float64(commonWords*2) / float64(totalWords)
}

func (s *MatchingService) levenshteinSimilarity(s1, s2 string) float64 {
    distance := s.levenshteinDistance(s1, s2)
    maxLen := math.Max(float64(len(s1)), float64(len(s2)))
    
    if maxLen == 0 {
        return 1.0
    }
    
    return 1.0 - (float64(distance) / maxLen)
}

// Simplified Jaro-Winkler similarity
func (s *MatchingService) jaroWinklerSimilarity(s1, s2 string) float64 {
    if s1 == s2 {
        return 1.0
    }
    
    // Simplified implementation
    commonChars := 0
    for i := 0; i < len(s1) && i < len(s2); i++ {
        if s1[i] == s2[i] {
            commonChars++
        }
    }
    
    if commonChars == 0 {
        return 0.0
    }
    
    maxLen := math.Max(float64(len(s1)), float64(len(s2)))
    return float64(commonChars) / maxLen
}

// Enhanced distance calculation
func (s *MatchingService) calculateDistance(from1, to1, from2, to2 string) float64 {
    fromDistance := s.stringDistance(from1, from2)
    toDistance := s.stringDistance(to1, to2)
    
    // Weight destination distance more heavily
    return (fromDistance * 0.4) + (toDistance * 0.6)
}

func (s *MatchingService) stringDistance(s1, s2 string) float64 {
    if s1 == s2 {
        return 0.0
    }
    
    maxLen := math.Max(float64(len(s1)), float64(len(s2)))
    if maxLen == 0 {
        return 0.0
    }
    
    distance := float64(s.levenshteinDistance(s1, s2))
    return (distance / maxLen) * 100
}

func (s *MatchingService) levenshteinDistance(s1, s2 string) int {
    if len(s1) == 0 {
        return len(s2)
    }
    if len(s2) == 0 {
        return len(s1)
    }
    
    matrix := make([][]int, len(s1)+1)
    for i := range matrix {
        matrix[i] = make([]int, len(s2)+1)
        matrix[i][0] = i
    }
    
    for j := 0; j <= len(s2); j++ {
        matrix[0][j] = j
    }
    
    for i := 1; i <= len(s1); i++ {
        for j := 1; j <= len(s2); j++ {
            cost := 0
            if s1[i-1] != s2[j-1] {
                cost = 1
            }
            
            matrix[i][j] = min(
                matrix[i-1][j]+1,      // deletion
                matrix[i][j-1]+1,      // insertion
                matrix[i-1][j-1]+cost, // substitution
            )
        }
    }
    
    return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
    if a < b && a < c {
        return a
    }
    if b < c {
        return b
    }
    return c
}

// Get top N matches
func (s *MatchingService) GetTopMatches(userID int, from, to string, maxDistance float64, limit int) ([]models.TripMatch, error) {
    allMatches, err := s.FindMatchingTrips(userID, from, to, maxDistance)
    if err != nil {
        return nil, err
    }
    
    if len(allMatches) <= limit {
        return allMatches, nil
    }
    
    return allMatches[:limit], nil
}
