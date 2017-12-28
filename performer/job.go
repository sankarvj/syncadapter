package performer

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var sql_create_jobs_table = `
	CREATE TABLE IF NOT EXISTS jobs(
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Job TEXT,
		Retry INTEGER DEFAULT 0,
		Status INTEGER DEFAULT 0,
		Updated DATETIME
	);
	`

var sql_job_add = `
	INSERT INTO jobs(
		Job,
		Retry,
		Status,
		Updated
	) values(?, ?, ?, CURRENT_TIMESTAMP)
	`

var sql_job_update = `
	UPDATE jobs set 
		Retry = ? 
		WHERE id = ?
	`

type Job struct {
	Id      int64
	Job     string
	Retry   int16
	Status  int16
	Updated time.Time
}

var db *sql.DB

//initDB must be called before excecuting delete operations
func InitJobDB(dbPath string) *sql.DB {
	var err error

	if db == nil {
		db, err = sql.Open("sqlite3", dbPath+"?mode=rwc")
		db.Exec("PRAGMA foreign_keys = ON;")
		panicError(err)
		if db == nil {
			panic("db nil")
		}
		createJobsTable(db)
	}
	return db
}

func createJobsTable(db *sql.DB) {
	// create account table if not exists
	_, err := db.Exec(sql_create_jobs_table)
	if err != nil {
		log.Println("error creating jobs table", err)
	}
}

func addJob(db *sql.DB, jobString string) {

	job := readAJob(db, jobString)
	if job.Id != 0 {
		if job.Retry < 2 {
			updateJobRetryCount(db, job.Id, job.Retry)
		} else {
			deleteJob(db, jobString)
			return
		}
	}

	stmt, err := db.Prepare(sql_job_add)
	defer stmt.Close()
	if err != nil {
		panicError(err)
		return
	}
	_, err = stmt.Exec(jobString, 0, 0)
	panicError(err)
}

func updateJobRetryCount(db *sql.DB, id int64, oldRetryCount int16) {
	stmt, err := db.Prepare(sql_job_update)
	panicError(err)
	_, err = stmt.Exec(oldRetryCount+1, id)
	panicError(err)
}

func deleteJob(db *sql.DB, jobString string) {
	stmt, err := db.Prepare("DELETE from jobs WHERE job = " + jobString)
	defer stmt.Close()
	if err != nil {
		panicError(err)
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		panicError(err)
		return
	}
}

func deleteAllJobs(db *sql.DB) {
	stmt, err := db.Prepare("DELETE from jobs")
	defer stmt.Close()
	if err != nil {
		panicError(err)
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		panicError(err)
		return
	}
}

func readAllJobs(db *sql.DB) []Job {
	sql_readall := `
	SELECT Id,Job,Retry,Status,Updated FROM jobs
	`
	rows, err := db.Query(sql_readall)
	defer closeRows(rows)
	if err != nil {
		log.Println("Error reading readAllJobs ", err)
		return
	}

	jobs := make([]Job, 0)
	for rows.Next() {
		job := &Job{}
		err = rows.Scan(&job.Id, &job.Job, &job.Retry, &job.Status, &job.Updated)
		if err != nil {
			log.Println("error scanning job row")
			continue
		}
		jobs = append(result, *job)
		log.Println("Job --- ", job)
	}
	return jobs
}

func readAJob(db *sql.DB, jobString string) Job {
	sql_readall := `
	SELECT Id,Job,Retry,Status,Updated FROM jobs
	WHERE Job = ` + jobString + ` LIMIT 1
	`
	rows, err := db.Query(sql_readall)
	defer closeRows(rows)
	if err != nil {
		log.Println("Error reading readAJob ", err)
		return
	}

	job := &Job{}
	for rows.Next() {
		err = rows.Scan(&job.Id, &job.Job, &job.Retry, &job.Status, &job.Updated)
		if err != nil {
			log.Println("Error scanning job row")
			continue
		}
		log.Println("Job --- ", job)
	}
	return *job
}

func panicError(err error) {
	if err != nil {
		log.Println("Don't forget to call initDB from syncAdapter ", err)
	}
}
