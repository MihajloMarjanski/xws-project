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
export class CompanyService {
  
  private _baseUrl = 'https://localhost:8600/';
  private _createCompany = this._baseUrl + 'company/create/';
  
  constructor(private _http: HttpClient) { }



  createCompany(company: Company): Observable<any> {
    const body=JSON.stringify(company);
    console.log(body)
    return this._http.post(this._createCompany + company.ownerUsername, body)
  }
  


  private handleError(err : HttpErrorResponse) {
    console.log(err.message);
    return Observable.throw(err.message);
    throw new Error('Method not implemented.');
  } 
}
