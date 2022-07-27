import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Credentials } from '../model/credentials';
import { Experience } from '../model/experience';
import { Interest } from '../model/interest';
import { JobOffer } from '../model/job-offer';
import { User } from '../model/user';

@Injectable({
  providedIn: 'root'
})
export class UserService {
 
  
  
  private _baseUrl = 'https://localhost:8000/';
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
  private _forgotPassword  = this._baseUrl + 'user/newPassword/';
  private _createOffer  = this._baseUrl + 'jobs/offer';
  private _apiKey  = this._baseUrl + 'user/apiKey/';
  private _searchUsers  = this._baseUrl + 'user/search/';
  private _searchOffers  = this._baseUrl + 'jobs/search/';
  private _blockUser  = this._baseUrl + 'user/block/';
  private _findAllBlocked  = this._baseUrl + 'user/blocked/';
  private _sendPasswordless = this._baseUrl + 'auth/sso/';
  private _send2factor = this._baseUrl + 'user/2factorAuth/pin/send';
  private _loginPaswordless = this._baseUrl + 'user/login/passwordless?token=';
  private _getRecommendedConnections = this._baseUrl + 'connection/users/';
  private _getUsersForIds = this._baseUrl + 'user/recommendedConnections/';
  
  constructor(private _http: HttpClient) { }
  

  getUsersForIds(id: any) {
    return this._http.get<any>(this._getUsersForIds + id)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  getRecommendedConnections(id: string | null) {
    return this._http.get<any>(this._getRecommendedConnections + id)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  loginPasswordless(token: any): Observable<any> {
    return this._http.post(this._loginPaswordless + token, {})
  } 

  sendPasswordless(username: string): Observable<any> {
    return this._http.get<any>(this._sendPasswordless + username)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  send2factor(credentials: Credentials): Observable<any> {
    const body=JSON.stringify(credentials);
    console.log(body)
    return this._http.post(this._send2factor, body)
  }

  findAllBlocked(id: string | null) {
    return this._http.get<any>(this._findAllBlocked + id)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError));
  }

  blockUser(logedUserId: string | null, blockedUserId: number) : Observable<any> {
    return this._http.put<any>(this._blockUser + logedUserId + "/" + blockedUserId, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }

  searchOffers(searchField: string) {
    return this._http.get<any>(this._searchOffers + searchField)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  searchUsers(username: string, id: string|null|number) {
    if (id == null)
      id = 0
    return this._http.get<any>(this._searchUsers + username + "/" + id)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  createOffer(offer: JobOffer): Observable<any>  {
    const body=JSON.stringify(offer);
    return this._http.post(this._createOffer, body)
  }

  getApiKey(username: string | null) {
    return this._http.get<any>(this._apiKey + username)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

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
