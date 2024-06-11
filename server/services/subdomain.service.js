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
    };
    export default subDomainService;