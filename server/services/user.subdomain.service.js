import pool from "../../dao/db-config.js";
import sql from "mssql";

const userSubDomainService = {
  insert: async (userEmail, subDomainName, callback) => {
    try {
      const result = await pool
        .request()
        .input("userEmail", sql.NVarChar, userEmail)
        .input("subDomainName", sql.NVarChar, subDomainName)
        .query(
          "INSERT INTO [studententuin].[dbo].[UserSubDomain] (userEmail, subDomainName) VALUES (@userEmail, @subDomainName)"
        );

      callback({status: 200, message: "user-subdomein succesvol toegevoegt"}, null);

    } catch (error) {
      callback(null, {status: 500, message: error.message});
    }
  }
};

export default userSubDomainService;
