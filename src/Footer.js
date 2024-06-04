import React from "react";
import { Link } from "react-router-dom";

export default function Footer() {
  return (
    <footer className="bg-neutral-warm">
      <div className="mx-auto flex max-w-full items-center justify-between p-6">
        <div className="flex flex-1 justify-start items-center">
          <Link to="/" className="-m-1.5 p-1.5">
            <img src="logo.png" className="w-16 h-16" />
          </Link>
        </div>
        <div className="flex flex-1 justify-end items-center gap-4 -px-2">
          <Link to="/about" className="text-lg font-semibold leading-6 text-black">
            over ons
          </Link>
          <Link to="/contact" className="text-lg font-semibold leading-6 text-black">
            contact
          </Link>
        </div>
      </div>
    </footer>
  );
}
