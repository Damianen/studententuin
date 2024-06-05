import pool from "../../dao/db-config.js";
import sql from "mssql";
import bcrypt from "bcrypt";

const requestService = {
  insert: async (email, password, subDomainName, productPackage, callback) => {

    const hashedPassword = await bcrypt.hash(password, 10);

    try {
      await pool
        .request()
        .input("email", sql.NVarChar, email)
        .input("password", sql.NVarChar, hashedPassword)
        .input("subDomainName", sql.NVarChar, subDomainName)
        .input("productPackage", sql.NVarChar, productPackage)
        .query(
          "INSERT INTO [studententuin].[dbo].[Request] (email, password, subDomainName, package) VALUES (@email, @password, @subDomainName, @productPackage)"
        );
    } catch (err) {
      callback(null, {message: err.message});
    }
  },
  getByEmail: async (userEmail, callback) => {
    try {
      const result = await pool
        .request()
        .input("userEmail", sql.NVarChar, userEmail)
        .query(
          "SELECT * FROM [studententuin].[dbo].[Request] WHERE email = @userEmail"
        );

      if (result.recordset.length !== 0) {
        callback({status: 200, data: result.recordset}, null);
      } else {
        callback(null, {status: 404, message: "User not found"});
      }
    } catch (error) {
      callback(null, {status: 500, message: error.message});
    }
  }
};

export default requestService;
