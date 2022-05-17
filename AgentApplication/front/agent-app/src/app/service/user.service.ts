import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Credentials } from '../model/credentials';
import { Observable, throwError } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';
import { CompanyOwner } from '../model/company-owner';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  private _baseUrl = 'http://localhost:8600/';
  private _login = this._baseUrl + 'auth/login';
  private _getAllUsernames = this._baseUrl + 'auth/getAllUsernames';
  private _submitRegistration  = this._baseUrl + 'company/owner/create';

  constructor(private _http: HttpClient) { }

  logIn(credentials: Credentials): Observable<any> {
    const body=JSON.stringify(credentials);
    console.log(body)
    return this._http.post(this._login, body,{responseType: 'text'})
  }

  getAllUsernames(): Observable<string[]> { 
    return this._http.get<string[]>(this._getAllUsernames)
                          .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                            catchError(this.handleError)); 
  }

  createCompanyOwner(owner: CompanyOwner) : Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
   });
    const body=JSON.stringify(owner);
    console.log(body)
    return this._http.post(this._submitRegistration, body, {'headers':headers})
  }

  private handleError(err : HttpErrorResponse) {
    console.log(err.message);
    return Observable.throw(err.message);
    throw new Error('Method not implemented.');
  } 
}
