import pool from "../../dao/db-config.js";
import sql from "mssql";
import bcrypt from "bcrypt";

const userService = {
  getByEmail: async (userEmail, callback) => {

    try {
      const result = await pool
        .request()
        .input("userEmail", sql.NVarChar, userEmail)
        .query(
          "SELECT * FROM [studententuin].[dbo].[User] WHERE Email = @userEmail"
        );

      if (result.recordset.length !== 0) {
        callback({status: 200, data: result.recordset}, null);
      } else {
        callback(null, {status: 404, message: "User not found"});
      }
    } catch (error) {
      callback(null, {status: 500, message: error.message});
    }
  },
  
  getUserByEmail: async (email) => {
    try {
      const request = pool.request();
      const result = await request
        .input("email", sql.NVarChar, email)
        .query(
          "SELECT * FROM [studententuin].[dbo].[User] WHERE Email = @email"
        );

      if (result.recordset.length > 0) {
        return result.recordset[0];
      } else {
        return null;
      }
    } catch (error) {
      console.error("Error:", error);
      throw error;
    }
  },

  register: async (email, password) => {
    try {
      const hashedPassword = await bcrypt.hash(password, 10);

      const result = await pool
        .request()
        .input("email", sql.NVarChar, email)
        .input("password", sql.NVarChar, hashedPassword)
        .query(
          "INSERT INTO [studententuin].[dbo].[User] (email, password) VALUES (@email, @password)"
        );

      return "User created successfully";

      console.log("User created successfully");
    } catch (error) {
      console.error("Error:", error);
      throw error;
    }
  },
};

export default userService;
