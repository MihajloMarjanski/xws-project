import { Comment } from "./comment";
import { JobPosition } from "./JobPosition";

export interface Company {
    id: number ;
    name: string ;
    info: string ;
    city: string ;
    country: string ;
    isApproved: boolean;
    ownerUsername: string;
    comments: Comment[];
    positions: JobPosition[];
}