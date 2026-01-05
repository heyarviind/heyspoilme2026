package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/config"
	"heyspoilme/internal/database"
	"heyspoilme/internal/models"
)

// Indian cities with coordinates for fake profiles
var indianCities = []struct {
	City      string
	State     string
	Latitude  float64
	Longitude float64
}{
	{"Mumbai", "Maharashtra", 19.0760, 72.8777},
	{"Delhi", "Delhi", 28.6139, 77.2090},
	{"Bangalore", "Karnataka", 12.9716, 77.5946},
	{"Hyderabad", "Telangana", 17.3850, 78.4867},
	{"Chennai", "Tamil Nadu", 13.0827, 80.2707},
	{"Kolkata", "West Bengal", 22.5726, 88.3639},
	{"Pune", "Maharashtra", 18.5204, 73.8567},
	{"Ahmedabad", "Gujarat", 23.0225, 72.5714},
	{"Jaipur", "Rajasthan", 26.9124, 75.7873},
	{"Lucknow", "Uttar Pradesh", 26.8467, 80.9462},
	{"Chandigarh", "Chandigarh", 30.7333, 76.7794},
	{"Gurgaon", "Haryana", 28.4595, 77.0266},
	{"Noida", "Uttar Pradesh", 28.5355, 77.3910},
	{"Indore", "Madhya Pradesh", 22.7196, 75.8577},
	{"Kochi", "Kerala", 9.9312, 76.2673},
	{"Goa", "Goa", 15.2993, 74.1240},
	{"Surat", "Gujarat", 21.1702, 72.8311},
	{"Nagpur", "Maharashtra", 21.1458, 79.0882},
	{"Coimbatore", "Tamil Nadu", 11.0168, 76.9558},
	{"Vizag", "Andhra Pradesh", 17.6868, 83.2185},
}

// Male display names
var maleNames = []string{
	"Arjun", "Rahul", "Vikram", "Aditya", "Rohan", "Karan", "Nikhil", "Akash", "Varun", "Siddharth",
	"Amit", "Raj", "Dev", "Aarav", "Ishaan", "Kabir", "Vihaan", "Reyansh", "Ayaan", "Krishna",
	"Shiv", "Arnav", "Dhruv", "Yash", "Aryan", "Pranav", "Shaurya", "Advait", "Rudra", "Atharv",
	"Vivaan", "Aayush", "Ansh", "Aarush", "Veer", "Sai", "Om", "Kiaan", "Rivaan", "Darsh",
	"Samarth", "Laksh", "Neil", "Virat", "Aaryan", "Kavish", "Shivansh", "Rehan", "Aadhya", "Ved",
}

// Female display names
var femaleNames = []string{
	"Priya", "Ananya", "Shreya", "Neha", "Pooja", "Riya", "Aisha", "Diya", "Isha", "Kavya",
	"Meera", "Nisha", "Simran", "Tanya", "Zara", "Aditi", "Bhavna", "Chitra", "Deepa", "Ekta",
	"Fatima", "Gauri", "Harini", "Ira", "Jhanvi", "Kiara", "Lavanya", "Mahi", "Natasha", "Ojaswi",
	"Pallavi", "Radhika", "Sakshi", "Tanvi", "Uma", "Vanshika", "Yashvi", "Zoya", "Aarna", "Bhumi",
	"Charvi", "Divya", "Esha", "Falguni", "Geet", "Hina", "Ishita", "Jahanvi", "Khushi", "Laxmi",
}

// Bio templates
var maleBios = []string{
	"Tech professional who loves to travel and explore new cultures. Looking for someone special.",
	"Entrepreneur building something exciting. Passionate about fitness and good food.",
	"Engineer by profession, adventurer by heart. Let's create amazing memories together.",
	"Finance guy who enjoys music, movies, and meaningful conversations.",
	"Doctor who believes in work-life balance. Love cooking and hiking on weekends.",
	"Startup founder with a love for innovation and adventure sports.",
	"Corporate professional who values honesty and genuine connections.",
	"Software developer who loves gaming, movies, and exploring cafes.",
	"Investment banker with a passion for travel and fine dining.",
	"Marketing professional who enjoys photography and outdoor activities.",
}

var femaleBios = []string{
	"Creative soul who loves art, music, and good coffee. Looking for genuine connections.",
	"Fashion enthusiast with a passion for travel and trying new cuisines.",
	"Corporate professional who enjoys reading, yoga, and weekend getaways.",
	"Doctor who believes in making a difference. Love music and dancing.",
	"Entrepreneur building my dreams. Passionate about fitness and wellness.",
	"Designer who finds beauty in everything. Let's explore life together.",
	"Lawyer with a love for travel and meaningful conversations.",
	"HR professional who enjoys cooking, movies, and quality time with loved ones.",
	"Content creator who loves storytelling and capturing beautiful moments.",
	"Finance professional who values ambition and genuine connections.",
}

var salaryRanges = []string{
	"5-10 LPA",
	"10-20 LPA",
	"20-50 LPA",
	"50+ LPA",
}

func main() {
	log.Println("Starting fake users seeding...")

	// Load config and connect to database
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	created := 0
	skipped := 0

	for i := 1; i <= 100; i++ {
		email := fmt.Sprintf("demos%d@aithreads.io", i)

		// Check if user already exists
		var existingID uuid.UUID
		err := db.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&existingID)
		if err == nil {
			log.Printf("User %s already exists, skipping", email)
			skipped++
			continue
		} else if err != sql.ErrNoRows {
			log.Printf("Error checking user %s: %v", email, err)
			continue
		}

		// Determine gender: odd = male, even = female
		var gender models.Gender
		var displayName string
		var bio string

		if i%2 == 1 {
			gender = models.GenderMale
			displayName = maleNames[rand.Intn(len(maleNames))]
			bio = maleBios[rand.Intn(len(maleBios))]
		} else {
			gender = models.GenderFemale
			displayName = femaleNames[rand.Intn(len(femaleNames))]
			bio = femaleBios[rand.Intn(len(femaleBios))]
		}

		// Random age between 21 and 45
		age := 21 + rand.Intn(25)

		// Random city
		city := indianCities[rand.Intn(len(indianCities))]

		// Random salary range for males
		salaryRange := ""
		if gender == models.GenderMale {
			salaryRange = salaryRanges[rand.Intn(len(salaryRanges))]
		}

		// Create user with email_verified = true
		userID := uuid.New()
		now := time.Now().UTC()

		_, err = db.Exec(`
			INSERT INTO users (id, email, email_verified, created_at, updated_at)
			VALUES ($1, $2, true, $3, $4)
		`, userID, email, now, now)
		if err != nil {
			log.Printf("Failed to create user %s: %v", email, err)
			continue
		}

		// Create fake profile
		profileID := uuid.New()
		var salaryRangeSQL interface{}
		if salaryRange != "" {
			salaryRangeSQL = salaryRange
		} else {
			salaryRangeSQL = nil
		}

		_, err = db.Exec(`
			INSERT INTO profiles (id, user_id, display_name, gender, age, bio, salary_range, city, state, latitude, longitude, is_complete, is_fake, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true, true, $12, $13)
		`, profileID, userID, displayName, gender, age, bio, salaryRangeSQL, city.City, city.State, city.Latitude, city.Longitude, now, now)
		if err != nil {
			log.Printf("Failed to create profile for %s: %v", email, err)
			// Rollback user creation
			db.Exec(`DELETE FROM users WHERE id = $1`, userID)
			continue
		}

		log.Printf("Created: %s - %s (%s, %d, %s)", email, displayName, gender, age, city.City)
		created++
	}

	log.Printf("Seeding complete! Created %d users, skipped %d existing", created, skipped)
}

