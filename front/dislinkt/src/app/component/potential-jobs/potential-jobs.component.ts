import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { JobOffer } from 'src/app/model/job-offer';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-potential-jobs',
  templateUrl: './potential-jobs.component.html',
  styleUrls: ['./potential-jobs.component.css']
})
export class PotentialJobsComponent implements OnInit {

  searchField: string = '';
  displayedColumns: string[] = ['companyName', 'jobPosition', 'jobInfo', 'qualifications'];
  errorMessage : string  = '';
  show : boolean = false
  offers: JobOffer[] = [
    {
      companyName: "DoBar",
      apiKey: "",
      jobInfo:  "Samo programiraj",
      jobPosition: "Programmer",
      qualifications: 'Python, Java, Matlab, Communication',
    },
    {
      companyName: "Klas",
      apiKey: "",
      jobInfo:  "Samo peci",
      jobPosition: "Baker",
      qualifications: 'Baker, Communication',
    },
    {
      companyName: "LaSorela",
      apiKey: "",
      jobInfo:  "Samo prodaj",
      jobPosition: "Salesman",
      qualifications: 'Communication',
    }
  ]


  constructor(public _userService: UserService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    //this.getRecommended()
  }

  getRecommended() {
    this._userService.getRecommendedConnections(this.searchField)
          .subscribe(data => {
              this.offers = data.offers
              this.show = false
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
    
  }

}
