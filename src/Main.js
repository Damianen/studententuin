import React from "react";

import Navigation from "./Navigation.js";
import Footer from "./Footer.js";

export default function Main({children, manage}) {
  return (
    <div className="min-h-screen flex justify-between flex-col">
      {manage === undefined ? <Navigation /> : <></>}
        {children}
      <Footer />
    </div>
  );
}
