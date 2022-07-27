import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Client } from 'src/app/model/client';
import { Comment } from 'src/app/model/comment';
import { Company } from 'src/app/model/company';
import { InterviewInformation } from 'src/app/model/InterviewInformation';
import { JobPosition } from 'src/app/model/JobPosition';
import { CompanyService } from 'src/app/service/company.service';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-forum',
  templateUrl: './forum.component.html',
  styleUrls: ['./forum.component.css']
})
export class ForumComponent implements OnInit {

  companies : Company[] = []
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
  currentClient! : Client
  comment: Comment = {
    client : this.currentClient,
    createdDate : new Date,
    id: 0,
    text: ''
  }
  interviewInfo: InterviewInformation = {
    client : this.currentClient,
    id: 0,
    info: ''
  }
  salary: number = 0

  constructor(public _userService: UserService, private _companyService: CompanyService, private router: Router, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
      this.getClient();
      this.getAllApprovedCompanies();
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

  addSalary() {
    this.selectedJob.avgSalary = this.salary
    this._companyService.addSalary(this.selectedJob)
        .subscribe(data =>  {
          this._snackBar.open('Salary updated.', 'Close', {duration: 3000}); 
          this.salary = 0
          },
        error => this.errorMessage = <any>error); 
  }

  createInterviewInfo() {
    this.interviewInfo.client = this.currentClient
    this._companyService.createInterviewInfo(this.selectedJob, this.interviewInfo)
        .subscribe(data =>  {
          this._snackBar.open('Interview information created.', 'Close', {duration: 3000});  
          this.interviewInfo.info = ''
          },
        error => this.errorMessage = <any>error);

  }

  createComment() {
    this.comment.client = this.currentClient
    this.comment.createdDate = new Date
    this._companyService.createComment(this.comment, this.selectedCompany.id)
        .subscribe(data =>  {
          this._snackBar.open('Comment created.', 'Close', {duration: 3000}); 
          this.comment.text = ''
          },
        error => this.errorMessage = <any>error);

  }
}
