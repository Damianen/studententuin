import pool from "../../dao/db-config.js";
import sql from "mssql";

const userService = {
  getByEmail: async (req, res, next) => {
    const userEmail = req.params.userEmail;

    try {
      const result = await pool
        .request()
        .input('userEmail', sql.NVarChar, userEmail)
        .query("SELECT * FROM [studententuin].[dbo].[User] WHERE Email = @userEmail");

      if (result.recordset) {
        res.send({
          status: 200,
          message: "This is a message",
          data: result.recordset,
        });
      }else{
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
  },
};

export default userService;
