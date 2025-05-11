
April 2nd - 3rd
After a night and a day, I finally connected to the google cloud database which I will be using
Have no idea why the code I'm using works above the other code i tried. 

Shockingly, this random medium article had the working code I needed:
https://medium.com/@DazWilkin/google-cloud-sql-6-ways-golang-a4aa497f3c67

And not:
https://adamtheautomator.com/cloud-sql/#Creating_a_Dedicated_Google_Cloud_SQL_Instance_User

Or the official: https://cloud.google.com/sql/docs/mysql/samples/cloud-sql-mysql-databasesql-connect-connector
+ chatgpt


My other attempts:

// user=postgres.snwmbzorkemzqnpwrwad password=[YOUR-PASSWORD] host=aws-0-us-west-1.pooler.supabase.com port=5432 dbname=postgres
	// Connect to DB
	/*cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST") + ":3306",
		DBName: os.Getenv("DB_NAME"),
	}*/

	//dsn := "user://postgres:Wc1ba51WjIdNR11e@db.snwmbzorkemzqnpwrwad.supabase.co:5432/postgres"
	//dsn := fmt.Sprintf("%s:%s@cloudsql-mysql(%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	// user=postgres password=[YOUR-PASSWORD] host=db.snwmbzorkemzqnpwrwad.supabase.co port=5432 dbname=postgres
	//db, err := sql.Open("cloudsql-mysql",
	//	fmt.Sprintf("%s:%s@cloudsql-mysql(%s)/%s",
	//		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST")))
	//db, err := sql.Open("mysql", cfg.FormatDSN())

	//db, err := sql.Open("mysql", dsn)

	// db, err := sql.Open("mysql", "user:password@tcp(instance-connection-name)/dbname")

	// db, err := connectWithConnectorIAMAuthN()
