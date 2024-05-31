import pool from "../../dao/db-config.js";

const userService = {
  getByEmail: async (req, res, next) => {
    const userEmail = req.params.userEmail;

    try {
      const result = await pool
        .request()
        .query("SELECT * FROM [studententuin].[dbo].[User]");

      if (result.recordset) {
        res.send({
          status: 200,
          message: "This is a message",
          data: result.recordset,
        });
      }
    } catch (error) {
      res.send(error);
    }
  },
};

export default userService;
