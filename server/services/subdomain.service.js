import pool from "../../dao/db-config.js";
import sql from "mssql";

const subDomainService = {
    getSubDomainByDomain: async (domain) => {
        try {
        const request = pool.request();
        const result = await request
            .input("domain", sql.NVarChar, domain)
            .query(
            "SELECT * FROM [studententuin].[dbo].[SubDomain] WHERE Domain = @domain"
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

    getRepo: async (subDomainName) => {
        try {
        const request = pool.request();
        const result = await request
            .input("name", sql.NVarChar, subDomainName)
            .query(
            "SELECT github FROM [studententuin].[dbo].[SubDomain] WHERE name = @name"
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
        getBranch: async (subDomainName) => {
            try {
            const request = pool.request();
            const result = await request
                .input("name", sql.NVarChar, subDomainName)
                .query(
                "SELECT githubBranch FROM [studententuin].[dbo].[SubDomain] WHERE name = @name"
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

    setBranch: async (subDomainName, branch) => {
        try {
        const request = pool.request();
        console.log(subDomainName, branch);
        const result = await request
            .input("name", sql.NVarChar, subDomainName)
            .input("branch", sql.NVarChar, branch)
            .query(
            "UPDATE [studententuin].[dbo].[SubDomain] SET githubBranch = @branch WHERE name = @name"
            );
    
            if (result.rowsAffected.length > 0) {
                console.log(result.rowsAffected[0]);
                return result.rowsAffected[0];
              } else {
                return null;
              }
            } catch (error) {
              console.error("Error:", error);
              throw error;
            }
        },
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

        insertNewRepo: async (subDomainName, github, branch, callback) => {
          try{
            const result = await pool
            .request()
            .input("subDomainName", sql.NVarChar, subDomainName)
            .input("github", sql.NVarChar, github)
            .input("branch", sql.NVarChar, branch)
            .query(
              "UPDATE [studententuin].[dbo].[Subdomain] SET github = @github, githubBranch = @branch WHERE name = @subDomainName");
            callback({status: 200, message: "Sub-domain succesvol toegevoegt"}, null);
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