import pool from "../../dao/db-config.js";
import sql from "mssql";

const userService = {
  getByEmail: async (req, res, next) => {
    const userEmail = req.params.userEmail;

    try {
      const result = await pool
        .request()
        .input("userEmail", sql.NVarChar, userEmail)
        .query(
          "SELECT * FROM [studententuin].[dbo].[User] WHERE Email = @userEmail"
        );

      if (result.recordset) {
        res.send({
          status: 200,
          message: "This is a message",
          data: result.recordset,
        });
      } else {
        res.send({
          status: 404,
          message: "User not found",
          data: {},
        });
      }
    } catch (error) {
      res.send({
        status: 500,
        message: "An error occurred",
        error: error.message,
      });
    }

    // Code voorbeeld om de call in de frontend te maken
    // fetch('http://localhost:3001/api/user/r@struijlaart.nl')
    // .then(response => {
    //   if (!response.ok) {
    //     throw new Error('Network response was not ok ' + response.statusText);
    //   }
    //   return response.json();
    // })
    // .then(data => {
    //   console.dir(data);
    // })
    // .catch(error => {
    //   console.dir(error);
    // });
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
  }

};

export default userService;
