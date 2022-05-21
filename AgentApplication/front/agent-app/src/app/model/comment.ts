import { Client } from "./client";

export interface Comment {
    id: number ;
    text: string ;
    createdDate: Date ;
    country: string ;
    isApproved: boolean;
    client: Client;
}