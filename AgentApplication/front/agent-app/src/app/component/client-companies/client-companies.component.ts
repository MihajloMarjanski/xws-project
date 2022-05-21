import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Client } from 'src/app/model/client';
import { Company } from 'src/app/model/company';
import { CompanyService } from 'src/app/service/company.service';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-client-companies',
  templateUrl: './client-companies.component.html',
  styleUrls: ['./client-companies.component.css']
})
export class ClientCompaniesComponent implements OnInit {

  companies : Company[] = []
  displayedColumns: string[] = ['name', 'info', 'owner'];
  displayedColumns1: string[] = ['text', 'date', 'client'];
  errorMessage : string  = '';
  selectedCompany: Company = {
    id: 0,
    name: "",
    info: "",
    isApproved: true,
    ownerUsername: "",
    city: "",
    country: "",
    comments: [],
    positions: []
  }
  role : string|null = localStorage.getItem('role');
  show: boolean = false;
  currentClient! : Client

  constructor(public _userService: UserService, private _companyService: CompanyService, private router: Router, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
    if(this.role == 'ROLE_CLIENT')
    {
      this.getClient();
      this.getAllApprovedCompanies();
    }   

  }
  

  getAllApprovedCompanies() {
    this._companyService.getAllApprovedCompanies()
        .subscribe(data =>  {
                    this.companies = data
                    //this.dataSource = new MatTableDataSource(this.weekendHouses);
                    },
                  error => this.errorMessage = <any>error);   
  }

  getClient() {
    this._userService.getClientByUsername(localStorage.getItem('username') || '')
              .subscribe(data => {
                this.currentClient = data
                console.log('Dobio: ', data)},
              error => this.errorMessage = <any>error);  
  }

  showInfo(company: Company) {
      this.selectedCompany = company
    this.show = true;
  }
}
