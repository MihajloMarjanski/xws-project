package com.example.agent.controller;

import com.example.agent.model.*;
import com.example.agent.model.dto.JobOffer;
import com.example.agent.model.dto.OwnerWithCompany;
import com.example.agent.model.dto.UserDto;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.ConfirmationTokenRepository;
import com.example.agent.service.ClientService;
import com.example.agent.service.CompanyService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.context.request.WebRequest;

import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;

@RestController
@RequestMapping(value = "/company")
@Slf4j
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
    public ResponseEntity<?> createCompanyOwner(@Valid @RequestBody UserDto companyOwner, BindingResult res, HttpServletRequest request) {
        if(res.hasErrors()){
            log.warn("Ip: {}, Fields for new company owner not valid!", request.getRemoteAddr());
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }
        return companyService.createCompanyOwner(companyOwner, request);
    }

    @PreAuthorize("hasAnyRole('POTENTIAL_COMPANY_OWNER', 'COMPANY_OWNER')")
    @PostMapping(path = "/owner/update")
    public ResponseEntity<?> updateCompanyOwner(@RequestBody OwnerWithCompany companyOwner , HttpServletRequest request) {
        return companyService.updateCompanyOwner(companyOwner, request);
    }

    @PreAuthorize("hasRole('POTENTIAL_COMPANY_OWNER')")
    @PostMapping(path = "/create/{ownerUsername}")
    public ResponseEntity<?> createCompany(@RequestBody Company company, @PathVariable String ownerUsername) {
        return companyService.sendCompanyRegistrationRequest(company, ownerUsername);
    }

    @PreAuthorize("hasRole('COMPANY_OWNER')")
    @GetMapping(path = "/owner/{id}")
    public ResponseEntity<?> getOwner(@PathVariable Integer id) {
        return companyService.getOwner(id);
    }

    @PreAuthorize("hasRole('CLIENT')")
    @PostMapping(path = "/comments/create/{companyId}")
    public ResponseEntity<?> createComment(@RequestBody Comment comment, @PathVariable Integer companyId) {
        return clientService.createComment(comment, companyId);
    }

    @PreAuthorize("hasRole('CLIENT')")
    @PostMapping(path = "/jobs/salary/update")
    public ResponseEntity<?> updateSalary(@RequestBody JobPosition jobSalary) {
        return clientService.updateSalary(jobSalary);
    }

    @PreAuthorize("hasRole('CLIENT')")
    @PostMapping(path = "/jobs/interview/{jobId}")
    public ResponseEntity<?> addInformation(@RequestBody InterviewInformation interviewInformation, @PathVariable Integer jobId) {
        return clientService.addInterviewInformation(interviewInformation, jobId);
    }

    @GetMapping(path = "/jobs/{companyId}")
    public ResponseEntity<?> allJobs(@PathVariable Integer companyId) {
        return companyService.getAllJobs(companyId);
    }

    @PreAuthorize("hasRole('COMPANY_OWNER')")
    @PostMapping(path = "/jobs/offer")
    public ResponseEntity<?> createJobOffer(@RequestBody JobOffer jobOffer) {
        return companyService.createJobOffer(jobOffer);
    }

    @GetMapping(path = "/owner/activate")
    public ResponseEntity<?> activateOwnerAccount(HttpServletRequest request, @RequestParam("token") String hashCode) {
        ConfirmationToken token = confirmationTokenRepository.findByConfirmationToken(hashCode);
        CompanyOwner companyOwner = token.getCompanyOwner();
        if (companyOwner == null || companyOwner.isActivated()) {
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }

        companyOwner.setActivated(true);
        companyOwnerRepository.save(companyOwner);
        log.info("Ip: {}, Username: {}, Client is successfully activated!",request.getRemoteAddr(), companyOwner.getUsername());


        HttpHeaders headers = new HttpHeaders();
        headers.add("Location", "https://localhost:4200");
        return new ResponseEntity<String>(headers, HttpStatus.OK);
    }

    @GetMapping(path = "/all")
    public ResponseEntity<?> allCompanies() {
        return companyService.getAll();
    }

    @GetMapping(path = "/all/approved")
    public ResponseEntity<?> allApprovedCompanies() {
        return companyService.getAllApproved();
    }

    @PreAuthorize("hasAnyRole('POTENTIAL_COMPANY_OWNER', 'COMPANY_OWNER')")
    @GetMapping(path = "/owner/username/{username}")
    public ResponseEntity<?> ownerByUsername(@PathVariable String username) {
        return companyService.getOwnerByUsername(username);
    }

    @GetMapping(path = "/{username}")
    public ResponseEntity<?> companyByOwnerUsername(@PathVariable String username) {
        return companyService.getByOwner(username);
    }

    @PreAuthorize("hasRole('COMPANY_OWNER')")
    @GetMapping(path = "/owner/apiKey/{username}/{password}")
    public ResponseEntity<?> apiKey(@PathVariable String username, @PathVariable String password) {
        return companyService.getApiKey(username, password);
    }

    @GetMapping(value = { "/jobs/search", "/jobs/search/{text}" })
    public ResponseEntity<?> searchOffers(@PathVariable(required = false) String text) {
        return companyService.searchOffers(text);
    }
}
