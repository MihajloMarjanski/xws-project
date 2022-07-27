import { Company } from "./company";

export interface CompanyOwner {
    id: number ;
    firstName: string ;
    lastName: string ;
    email: string ;
    username: string;
    password: string ;
    company: Company;
}