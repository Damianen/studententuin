import NextAuth, { getServerSession } from "next-auth";

import { authOptions } from "./auth-options";

const handler = NextAuth(authOptions);

export const authHandlers = handler;
export const auth = () => getServerSession(authOptions);
export const signIn = handler.signIn;
export const signOut = handler.signOut;
