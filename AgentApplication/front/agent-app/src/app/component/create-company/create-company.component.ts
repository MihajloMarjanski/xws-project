import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Company } from 'src/app/model/company';
import { JobPosition } from 'src/app/model/JobPosition';
import { CompanyService } from 'src/app/service/company.service';

@Component({
  selector: 'app-create-company',
  templateUrl: './create-company.component.html',
  styleUrls: ['./create-company.component.css']
})
export class CreateCompanyComponent implements OnInit {

  company: Company = {
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

  newJob: JobPosition = {
      id: 0,
      name: "",
      avgSalary: 0,
      interviewInformations: []
  }

  errorMessage : string  = '';
  selectedPositions : JobPosition[] = []
  username: string|null = localStorage.getItem('username');

  constructor(private _companyService: CompanyService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    if (this.username != null)
      this.company.ownerUsername = this.username
  }


  addNewJob() {
    let newJob  = Object.assign({},this.newJob);
    this.company.positions.push(newJob);
    this.newJob.name = '';
  }

  createCompany()
  {
    this._companyService.createCompany(this.company)
      .subscribe(data => {
        console.log('Dobio: ', data) 
        this._snackBar.open('Company registration request has been successfully created!', 'Close', {duration: 5000}); 
        },
        error => {
          this.errorMessage = <any>error
          this._snackBar.open(error.error, 'Close', {duration: 4000});
        }); 

    console.log(this.company);
    this.router.navigateByUrl("/");
  }

}
