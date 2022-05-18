package com.example.agent.controller;

import com.example.agent.model.*;
import com.example.agent.model.dto.JobOffer;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.ConfirmationTokenRepository;
import com.example.agent.service.ClientService;
import com.example.agent.service.CompanyService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.context.request.WebRequest;

@RestController
@RequestMapping(value = "/company")
public class CompanyController {

    @Autowired
    private CompanyService companyService;
    @Autowired
    ClientService clientService;
    @Autowired
    ConfirmationTokenRepository confirmationTokenRepository;
    @Autowired
    private CompanyOwnerRepository companyOwnerRepository;


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

    @GetMapping(path = "/owner/activate")
    public ResponseEntity<?> activateOwnerAccount(WebRequest request, @RequestParam("token") String hashCode) {
        ConfirmationToken token = confirmationTokenRepository.findByConfirmationToken(hashCode);
        CompanyOwner companyOwner = token.getCompanyOwner();
        if (companyOwner == null || companyOwner.isActivated()) {
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }

        companyOwner.setActivated(true);
        companyOwnerRepository.save(companyOwner);

        HttpHeaders headers = new HttpHeaders();
        headers.add("Location", "http://localhost:4200");
        return new ResponseEntity<String>(headers, HttpStatus.OK);
    }
}
