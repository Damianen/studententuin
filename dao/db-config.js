import sql from "mssql";
import "dotenv/config";

const poolConfig = {
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  server: process.env.DB_HOST,
  port: Number(process.env.DB_PORT),
  database: process.env.DB_DATABASE,
  options: {
    encrypt: true,
    trustServerCertificate: true,
  },
};

const pool = new sql.ConnectionPool(poolConfig);

async function initializePool() {
  try {
    await pool.connect();
    console.log("Database Connected");

    pool.on("error", (err) => {
      console.log("Unexpected error on idle connection pool, error: " + err);
    });
  } catch (error) {
    console.log("Error connecting database: " + error);
  }
}

initializePool();

export default pool;
