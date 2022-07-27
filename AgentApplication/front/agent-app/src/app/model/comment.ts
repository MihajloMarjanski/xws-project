import { Client } from "./client";

export interface Comment {
    id: number ;
    text: string ;
    createdDate: Date ;
    client: Client;
}