import { InterviewInformation } from "./InterviewInformation";

export interface JobPosition {
    id: number ;
    name: string ;
    avgSalary: number ;
    interviewInformations: InterviewInformation[];
}