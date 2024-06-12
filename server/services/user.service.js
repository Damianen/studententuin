import pool from "../../dao/db-config.js";
import sql from "mssql";
import bcrypt from "bcryptjs"

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
          "SELECT * FROM [studententuin].[dbo].[User] WHERE email = @email"
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

  getUserByEmailFromSession: async (req) => {
    try {
      // Haal de e-mail van de gebruiker uit de sessie
      const userEmail = req.session.user.email;

      // Roep de getUserByEmail-functie aan om de gebruiker op te halen
      const request = pool.request();
      const result = await request
        .input("userEmail", sql.NVarChar, userEmail)
        .query(
          "SELECT * FROM [studententuin].[dbo].[UserSubDomain] WHERE userEmail = @userEmail"
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

  createUser: async (email, password, userPackage) => {
    try {
      const request = pool.request();
      const result = await request
        .input("email", sql.NVarChar, email)
        .input("password", sql.NVarChar, password)
        .input("userPackage", sql.NVarChar, userPackage)
        .query(
          "INSERT INTO [studententuin].[dbo].[User] (Email, Password, package) VALUES (@email, @password, @userPackage)"
        );

      return {
        email: email,
        userPackage: userPackage,
      };
    } catch (error) {
      console.error("Error:", error);
      throw error;
    }
  },

  getUserPackage: async (req) => {
    try {
      // Haal de e-mail van de gebruiker uit de sessie
      const userEmail = req.session.user.email;

      // Roep de getUserByEmail-functie aan om de gebruiker op te halen
      const request = pool.request();
      const result = await request
        .input("userEmail", sql.NVarChar, userEmail)
        .query(
          "SELECT package FROM [studententuin].[dbo].[User] WHERE email = @userEmail"
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

  getSubdomainByUser: async (req) => {
    try {
      // Haal de e-mail van de gebruiker uit de sessie
      const userEmail = req.session.user.email;

      // Roep de getUserByEmail-functie aan om de gebruiker op te halen
      const request = pool.request();
      const result = await request
        .input("userEmail", sql.NVarChar, userEmail)
        .query(
          "SELECT SubDomainName FROM [studententuin].[dbo].[UserSubDomain] WHERE userEmail = @userEmail"
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
  insert: async (email, password, productPackage, callback) => {
    
    const hashedPassword = await bcrypt.hash(password, 10);
    
    try {

      const result = await pool
        .request()
        .input("email", sql.NVarChar, email)
        .input("password", sql.NVarChar, hashedPassword)
        .input("package", sql.NVarChar, productPackage)
        .query(
          "INSERT INTO [studententuin].[dbo].[User] (email, password, package) VALUES (@email, @password, @package)"
        );

      callback({status: 200, message: "User succesvol toegevoegt"}, null);

    } catch (err) {
      callback(null, {status: 500, message: err.message});
    }
  }

};

export default userService;
