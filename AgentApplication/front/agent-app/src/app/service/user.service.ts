import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Credentials } from '../model/credentials';
import { Observable, throwError } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';
import { Client } from '../model/client';
import { CompanyOwner } from '../model/company-owner';
import { Company } from '../model/company';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  

  private _baseUrl = 'https://localhost:8600/';
  private _login = this._baseUrl + 'auth/login';
  private _getAllUsernames = this._baseUrl + 'auth/getAllUsernames';
  private _getCompanies = this._baseUrl + 'company/all';
  private _blackList = this._baseUrl + 'auth/password/blackList/';
  private _submitRegistrationOwner  = this._baseUrl + 'company/owner/create';
  private _submitRegistrationClient  = this._baseUrl + 'clients/create';
  private _clientByUsername  = this._baseUrl + 'clients/username/';
  private _adminByUsername  = this._baseUrl + 'admin/username/';
  private _ownerByUsername  = this._baseUrl + 'company/owner/username/';
  private _approveCompany  = this._baseUrl + 'admin/approve/company/';
  private _editOwner  = this._baseUrl + 'company/owner/update';
  private _editClient  = this._baseUrl + 'clients/update';
  private _editAdmin  = this._baseUrl + 'admin/update';
  private _forgotPassword  = this._baseUrl + 'clients/newPassword/';

  constructor(private _http: HttpClient) { }


  forgotPassword(username: string) : Observable<any> {
    return this._http.put<any>(this._forgotPassword + username, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }

  getClientByUsername(username: string): Observable<Client> {
    return this._http.get<Client>(this._clientByUsername + username)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  getAdminByUsername(username: string): Observable<Client> {
    return this._http.get<Client>(this._adminByUsername + username)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  getCompanyOwnerByUsername(username: string): Observable<CompanyOwner> {
    return this._http.get<CompanyOwner>(this._ownerByUsername + username)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  getCompanies(): Observable<Company[]> {
    return this._http.get<Company[]>(this._getCompanies)
                          .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                            catchError(this.handleError)); 
  }

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

  checkBlackListPass(pass: string): Observable<any> { 
    return this._http.get<any>(this._blackList + pass, {responseType: 'text' as 'json'})
                          .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                            catchError(this.handleError)); 
  }

  createCompanyOwner(owner: CompanyOwner) : Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
   });
    const body=JSON.stringify(owner);
    console.log(body)
    return this._http.post(this._submitRegistrationOwner, body, {'headers':headers})
  }

  createClient(client: Client) : Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
   });
    const body=JSON.stringify(client);
    console.log(body)
    return this._http.post(this._submitRegistrationClient, body, {'headers':headers})
  }

  approveCompany(id: number) : Observable<any> {
    return this._http.put<any>(this._approveCompany + id, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }

  editOwner(owner: CompanyOwner) : Observable<any> {
    const body=JSON.stringify(owner);
    console.log(body)
    return this._http.post(this._editOwner, body)
  }

  editClient(owner: Client) : Observable<any> {
    const body=JSON.stringify(owner);
    console.log(body)
    return this._http.post(this._editClient, body)
  }

  editAdmin(owner: Client) : Observable<any> {
    const body=JSON.stringify(owner);
    console.log(body)
    return this._http.post(this._editAdmin, body)
  }
  


  private handleError(err : HttpErrorResponse) {
    console.log(err.message);
    return Observable.throw(err.message);
    throw new Error('Method not implemented.');
  } 
}
