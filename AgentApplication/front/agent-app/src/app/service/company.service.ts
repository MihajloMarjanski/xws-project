import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Credentials } from '../model/credentials';
import { Observable, throwError } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';
import { Client } from '../model/client';
import { CompanyOwner } from '../model/company-owner';
import { Company } from '../model/company';
import { JobPosition } from '../model/JobPosition';
import { InterviewInformation } from '../model/InterviewInformation';
import { Comment } from '../model/comment';
import { JobOffer } from '../model/job-offer';

@Injectable({
  providedIn: 'root'
})
export class CompanyService {
  
  private _baseUrl = 'https://localhost:8600/';
  private _createCompany = this._baseUrl + 'company/create/';
  private _getCompanies = this._baseUrl + 'company/all/approved';
  private _createComment = this._baseUrl + 'company/comments/create/';
  private _createInterviewInfo = this._baseUrl + 'company/jobs/interview/';
  private _addSalary = this._baseUrl + 'company/jobs/salary/update';
  private _createOffer = this._baseUrl + 'company/jobs/offer';
  
  constructor(private _http: HttpClient) { }
  

  createOffer(offer: JobOffer): Observable<any>  {
    const body=JSON.stringify(offer);
    return this._http.post(this._createOffer, body)
  }

  createComment(comment: Comment, id: number): Observable<any>  {
    const body=JSON.stringify(comment);
    return this._http.post(this._createComment + id, body)
  }
  createInterviewInfo(selectedJob: JobPosition, interviewInfo: InterviewInformation)  : Observable<any> {
    const body=JSON.stringify(interviewInfo);
    return this._http.post(this._createInterviewInfo + selectedJob.id, body)
  }
  addSalary(selectedJob: JobPosition) : Observable<any>  {
    const body=JSON.stringify(selectedJob);
    return this._http.post(this._addSalary, body)
  }

  createCompany(company: Company): Observable<any> {
    const body=JSON.stringify(company);
    console.log(body)
    return this._http.post(this._createCompany + company.ownerUsername, body)
  }

  getAllApprovedCompanies(): Observable<Company[]> {
    return this._http.get<Company[]>(this._getCompanies)
                          .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                            catchError(this.handleError)); 
  }
  


  private handleError(err : HttpErrorResponse) {
    console.log(err.message);
    return Observable.throw(err.message);
    throw new Error('Method not implemented.');
  } 
}
