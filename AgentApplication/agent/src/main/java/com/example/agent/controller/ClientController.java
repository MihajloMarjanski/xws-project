package com.example.agent.controller;

import com.example.agent.model.*;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.ConfirmationTokenRepository;
import com.example.agent.service.ClientService;
import com.example.agent.service.CompanyService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.context.request.WebRequest;

import java.time.Period;
import java.util.Date;

@RestController
@RequestMapping(value = "/clients")
public class ClientController {

    @Autowired
    ClientService clientService;
    @Autowired
    ConfirmationTokenRepository confirmationTokenRepository;
    @Autowired
    CompanyService companyService;


    @PostMapping(path = "/create")
    public ResponseEntity<?> create(@RequestBody Client client) {
        return clientService.create(client);
    }

    @PreAuthorize("hasRole('CLIENT')")
    @PostMapping(path = "/update")
    public ResponseEntity<?> updateClient(@RequestBody Client client) {
        return clientService.updateClient(client);
    }

    @GetMapping(path = "/activate")
    public ResponseEntity<?> activateClientAccount(WebRequest request, @RequestParam("token") String hashCode) {
        ConfirmationToken token = confirmationTokenRepository.findByConfirmationToken(hashCode);
        Long secs = (token.getCreatedDate().getTime() - new Date().getTime())/1000;
        Client verificationClient = token.getClient();
        if (verificationClient == null || verificationClient.isActivated() || secs > 3600 ) {
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }

        verificationClient.setActivated(true);
        clientService.save(verificationClient);

        HttpHeaders headers = new HttpHeaders();
        headers.add("Location", "https://localhost:4200");
        return new ResponseEntity<String>(headers, HttpStatus.OK);
    }

    @PutMapping(path = "/newPassword/{email}")
    public ResponseEntity<?> sendNewPassword(@PathVariable String email) {
        if(clientService.findByEmail(email) != null)
            return clientService.sendNewPassword(clientService.findByEmail(email));
        if(companyService.findByOwnerEmail(email) != null)
            return companyService.sendNewPassword(companyService.findByOwnerEmail(email));
        return new ResponseEntity<String>(HttpStatus.BAD_REQUEST);
    }

    @PreAuthorize("hasRole('CLIENT')")
    @GetMapping(path = "/username/{username}")
    public ResponseEntity<?> clientByUsername(@PathVariable String username) {
        return clientService.getByUsername(username);
    }
}
