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
};

export default userService;
