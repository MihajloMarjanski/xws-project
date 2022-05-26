import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Company } from 'src/app/model/company';
import { Credentials } from 'src/app/model/credentials';
import { JobOffer } from 'src/app/model/job-offer';
import { JobPosition } from 'src/app/model/JobPosition';
import { CompanyService } from 'src/app/service/company.service';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-job-offer',
  templateUrl: './job-offer.component.html',
  styleUrls: ['./job-offer.component.css']
})
export class JobOfferComponent implements OnInit {

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
  selectedJob : JobPosition = {
    id: 0,
    name: "",
    avgSalary: 0,
    interviewInformations: []
  }
  offer: JobOffer = {
    qualifications: '',
    companyName: '',
    jobInfo: '',
    jobPosition: '',
    apiKey: ''
  }
  credentials: Credentials = {
    password : '',
    username : '',
    pin : ''
  }


  constructor(public _userService: UserService, private _companyService: CompanyService, private router: Router, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
    this.getCompany();
  }



  getCompany() {
    this._userService.getCompanyByOwnerUsername(localStorage.getItem('username') || '')
              .subscribe(data => {
                this.selectedCompany = data
                console.log('Dobio: ', data)},
              error => this.errorMessage = <any>error);  
  }

  createOffer() {
    this._userService.getApiKey(this.credentials.username, this.credentials.password) 
        .subscribe(data => {
          this.offer.apiKey = data
          this.offer.companyName = this.selectedCompany.name
          this.offer.jobPosition = this.selectedJob.name
          
          this._companyService.createOffer(this.offer)
              .subscribe(data => {
                this.router.navigateByUrl('/').then(() => {
                  this._snackBar.open('Job offer has been successfully created.', 'Close', {duration: 5000});
                });
                console.log('Dobio: ', data)},
              error => this.errorMessage = <any>error);  

          console.log('Dobio: ', data)},
        error => this.errorMessage = <any>error);

  }

}
