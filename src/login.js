import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";

const Login = () => {
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState(""); // Stukje state om foutmelding bij te houden

  useEffect(() => {
    const sessionCookie = Cookies.get("connect.sid");
    if (sessionCookie) {
      navigate("/manage");
    }
  }, [navigate]);

  const handleSubmit = async (event) => {
    event.preventDefault();
    const email = event.target.email.value;
    const password = event.target.password.value;

    const response = await fetch("/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (response.ok && data.redirectUrl) {
      // Zet de sessie cookie
      Cookies.set("connect.sid", data.sessionId, { expires: 1 }); // Verloopt na 1 dag
      navigate(data.redirectUrl);
    } else {
      setErrorMessage(data.message); // Zet de foutmelding in de state
    }
  };

  return (
    <div className="flex justify-center items-center h-screen">
      <div className="justify-center items-center">
        <h1 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900 p-5">
          Login
        </h1>
        <div className="lg:p-28 md:p-36 sm:20 p-4 w-full shadow-md border border-gray-200 rounded-lg">
          {/* Toon de foutmelding als deze aanwezig is */}
          {errorMessage && <p className="text-red-500">{errorMessage}</p>}
          <form onSubmit={handleSubmit}>
            {/* Email Input */}
            <div className="mb-4">
              <label
                htmlFor="email"
                className="block text-sm font-medium leading-6 text-gray-900"
              >
                Email
              </label>
              <input
                type="email"
                id="email"
                name="email"
                className="w-full border border-gray-300 rounded-md py-2 px-3 focus:outline-none focus:border-green-400"
                autoComplete="off"
              ></input>
            </div>
            {/* Password Input */}
            <div className="mb-6">
              <label
                htmlFor="password"
                className="block text-sm font-medium leading-6 text-gray-900"
              >
                Password
              </label>
              <input
                type="password"
                id="password"
                name="password"
                className="w-full border border-gray-300 rounded-md py-2 px-3 focus:outline-none focus:border-green-400"
              ></input>
            </div>
            {/* Remember me and Forgot password Link */}
            <div className="flex items-center">
              <input
                type="checkbox"
                id="remember"
                name="remember"
                className="text-green-400"
              ></input>
              <label htmlFor="remember" className="text-gray-600 ml-2">
                Onthoud mij
              </label>
            </div>
            <div className="text-primary-green">
              <Link to="wachtwoord" className="hover:underline">
                Wachtwoord vergeten?
              </Link>
            </div>
            <div className="mb-6 text-primary-green">
              <a href="#" className="hover:underline text-primary-green h-100%">
                Sign up Here
              </a>
            </div>
            {/* Sign up Link */}
            <div className="flex justify-between items-center">
              <button
                type="submit"
                className="rounded-full bg-gray-500 px-5 py-2 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-gray-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-600"
              >
                Back
              </button>
              {/* Login Button */}
              <button
                type="submit"
                className="rounded-full bg-primary-green px-5 py-2 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-house-green focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
              >
                Login
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Login;
