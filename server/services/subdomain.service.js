import pool from "../../dao/db-config.js";
import sql from "mssql";

const subDomainService = {
  getByName: async (subDomainName, callback) => {
    try {
      const result = await pool
        .request()
        .input("subDomainName", sql.NVarChar, subDomainName)
        .query(
          "SELECT * FROM [studententuin].[dbo].[SubDomain] WHERE name = @subDomainName"
        );

      if (result.recordset.length !== 0) {
        callback({status: 200, data: result.recordset}, null);
      } else {
        callback(null, {status: 404, message: "Sub-domain not found"});
      }
    } catch (error) {
      callback(null, {status: 500, message: error.message});
    }
  },

  insertDomainName: async (subDomainName, callback) => {
    try {
      const result = await pool
        .request()
        .input("subDomainName", sql.NVarChar, subDomainName)
        .query(
          "INSERT INTO [studententuin].[dbo].[SubDomain] (name) VALUES (@subDomainName)"
        );

      callback({status: 200, message: "Sub-domain succesvol toegevoegt"}, null);

    } catch (error) {
      callback(null, {status: 500, message: error.message});
    }
  }
};

export default subDomainService;
