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
  private _approveCompany = this._baseUrl + 'admin/approve/company/';

  constructor(private _http: HttpClient) { }
}
