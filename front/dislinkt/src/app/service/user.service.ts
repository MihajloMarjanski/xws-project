import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Credentials } from '../model/credentials';
import { Observable, throwError } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';
import { User } from '../model/user';
import { Experience } from '../model/experience';
import { Interest } from '../model/interest';
/* import { Client } from '../model/client';
import { CompanyOwner } from '../model/company-owner';
import { Company } from '../model/company'; */

@Injectable({
  providedIn: 'root'
})
export class UserService {
  
  private _baseUrl = 'http://localhost:8000/';
  private _login = this._baseUrl + 'user/login';
  private _getAllUsernames = this._baseUrl + 'auth/getAllUsernames';
  private _removeExperience = this._baseUrl + 'user/experience/';
  private _blackList = this._baseUrl + 'auth/password/blackList/';
  private _addExperience  = this._baseUrl + 'user/experience';
  private _submitRegistration  = this._baseUrl + 'user';
  private _getUserByUsername  = this._baseUrl + 'user/username/';
  private _addInterest  = this._baseUrl + 'user/interest';
  private _removeInterest  = this._baseUrl + 'user/interest/';
  private _editUser  = this._baseUrl + 'user';
  private _forgotPassword  = this._baseUrl + 'clients/newPassword/';
  private _companyByOwnerUsername  = this._baseUrl + 'company/';
  private _apiKey  = this._baseUrl + 'company/owner/apiKey/';
  
  constructor(private _http: HttpClient) { }
  

  getUserByUsername(username: string): Observable<any> {
    return this._http.get<any>(this._getUserByUsername + username)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  addInterest(newInterest: Interest) {
    const body=JSON.stringify(newInterest);
    console.log(body)
    return this._http.post(this._addInterest, body)
  }

  addExperience(newExperience: Experience) {
    const body=JSON.stringify(newExperience);
    console.log(body)
    return this._http.post(this._addExperience, body)
  }

  removeInterest(id: number) : Observable<any> {
    return this._http.delete<any>(this._removeInterest + id, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }

  removeExperience(id: number) : Observable<any> {
    return this._http.delete<any>(this._removeExperience + id, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }

  forgotPassword(username: string) : Observable<any> {
    return this._http.put<any>(this._forgotPassword + username, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }

  logIn(credentials: Credentials): Observable<any> {
    const body=JSON.stringify(credentials);
    console.log(body)
    return this._http.post(this._login, body)
  }

  getAllUsernames(): Observable<string[]> { 
    return this._http.get<string[]>(this._getAllUsernames)
                          .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                            catchError(this.handleError)); 
  }

  checkBlackListPass(pass: string): Observable<any> { 
    return this._http.get<any>(this._blackList + pass, {responseType: 'text' as 'json'})
                          .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                            catchError(this.handleError)); 
  }

  createUser(client: User) : Observable<any> {
    client.phone = client.phone.toString()
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
   });
    const body=JSON.stringify(client);
    console.log(body)
    return this._http.post(this._submitRegistration, body, {'headers':headers})
  }
/*
  approveCompany(id: number) : Observable<any> {
    return this._http.put<any>(this._approveCompany + id, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }
*/

  editUser(user: User) : Observable<any> {
    const body=JSON.stringify(user);
    console.log(body)
    return this._http.put<any>(this._editUser, body)
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  } 
  


  private handleError(err : HttpErrorResponse) {
    console.log(err.message);
    return Observable.throw(err.message);
    throw new Error('Method not implemented.');
  } 
}
