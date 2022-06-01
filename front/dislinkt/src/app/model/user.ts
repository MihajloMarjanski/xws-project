import { Experience } from "./experience";
import { Interest } from "./interest";

export interface User {
    id: number ;
    name: string ;
    gender: string ;
    email: string ;
    username: string;
    password: string ;
    phone: string ;
    date: Date|string| null;
    biography: string ;
    experiences: Experience[]
    interests: Interest[]
    isPrivate: boolean
}