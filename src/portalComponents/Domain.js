import React, { useState, useEffect } from "react";
import axios from "axios";

function Domain() {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await axios.get(
                    "/api/getUserByEmailFromSession"
                );
                console.log("API response:", response.data);
                setUser(response.data);
            } catch (error) {
                console.error("Error fetching user:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchUser();
    }, []);

    const handleClick = () => {
        if (!loading && user && user.subDomainName) {
            window.location.href = `https://${user.subDomainName}.studententuin.nl`;
        } else {
            alert("No user data available");
        }
    };

    return (
        <div>
            <h1 className="mb-10">
                Dit is de domein pagina, hier kan je een preview zien van je
                gedeployde website
            </h1>
            <div
                style={{
                    position: "relative",
                    cursor: "pointer",
                    width: "60vw",
                    height: "60vh",
                    overflow: "hidden",
                }}
                onClick={handleClick}
            >
                {loading ? (
                    <p>Laden...</p>
                ) : user && user.subDomainName ? (
                    <iframe
                        src={`https://${user.subDomainName}.studententuin.nl`}
                        style={{
                            width: "166.67%",
                            height: "166.67%",
                            transform: "scale(0.6)",
                            transformOrigin: "0 0",
                            border: "none",
                            pointerEvents: "none",
                        }}
                        title="Studententuin"
                    ></iframe>
                ) : (
                    <p>Geen website beschikbaar</p>
                )}
                <div
                    style={{
                        position: "absolute",
                        top: 0,
                        left: 0,
                        width: "100%",
                        height: "100%",
                    }}
                ></div>
            </div>
        </div>
    );
}

export default Domain;
