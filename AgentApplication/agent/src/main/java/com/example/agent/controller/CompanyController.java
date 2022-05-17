package com.example.agent.controller;

import com.example.agent.model.*;
import com.example.agent.model.dto.JobOffer;
import com.example.agent.service.ClientService;
import com.example.agent.service.CompanyService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(value = "/company")
public class CompanyController {

    @Autowired
    private CompanyService companyService;
    @Autowired
    ClientService clientService;


    @PostMapping(path = "/owner/create")
    public ResponseEntity<?> createCompanyOwner(@RequestBody CompanyOwner companyOwner) {
        return companyService.saveCompanyOwner(companyOwner);
    }

    @PostMapping(path = "/create")
    public ResponseEntity<?> createCompany(@RequestBody Company company) {
        return companyService.sendCompanyRegistrationRequest(company);
    }

    @GetMapping(path = "/owner/{id}")
    public ResponseEntity<?> getOwner(@PathVariable Integer id) {
        return companyService.getOwner(id);
    }


    @PostMapping(path = "/comments/create")
    public ResponseEntity<?> createComment(@RequestBody Comment comment) {
        return clientService.createComment(comment);
    }

    @PostMapping(path = "/jobs/salary/update")
    public ResponseEntity<?> updateSalary(@RequestBody JobPosition jobSalary) {
        return clientService.updateSalary(jobSalary);
    }

    @PostMapping(path = "/jobs/interview")
    public ResponseEntity<?> addInformation(@RequestBody InterviewInformation interviewInformation) {
        return clientService.addInterviewInformation(interviewInformation);
    }

    @GetMapping(path = "/jobs/{companyId}")
    public ResponseEntity<?> allJobs(@PathVariable Integer companyId) {
        return companyService.getAllJobs(companyId);
    }

    @PostMapping(path = "/jobs/offer")
    public ResponseEntity<?> createJobOffer(@RequestBody JobOffer jobOffer) {
        return companyService.createJobOffer(jobOffer);
    }
}
