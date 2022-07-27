package com.example.agent.controller;

import com.example.agent.model.*;
import com.example.agent.model.dto.UserCredentials;
import com.example.agent.repository.ConfirmationTokenRepository;
import com.example.agent.security.tokenUtils.JwtTokenUtils;
import com.example.agent.service.AdminService;
import com.example.agent.service.ClientService;
import com.example.agent.service.CompanyService;
import com.example.agent.service.EmailService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.context.request.WebRequest;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.net.URISyntaxException;
import java.util.Calendar;
import java.util.Date;
import java.util.HashSet;
import java.util.Set;


//Kontroler zaduzen za autentifikaciju korisnika
@Slf4j
@RestController
@RequestMapping(value = "/auth", produces = MediaType.APPLICATION_JSON_VALUE)
public class AuthenticationController {

    @Autowired
    ConfirmationTokenRepository confirmationTokenRepository;
    @Autowired
    EmailService emailService;
    private final JwtTokenUtils tokenUtils;
    private final AuthenticationManager authenticationManager;
    private final AdminService adminService;
    private final CompanyService companyService;
    private final ClientService clientService;

    @Autowired
    public AuthenticationController(AdminService adminService, CompanyService companyService, ClientService clientService,
                                    AuthenticationManager authenticationManager, JwtTokenUtils tokenUtils) {
        this.tokenUtils = tokenUtils;
        this.authenticationManager = authenticationManager;
        this.adminService = adminService;
        this.companyService = companyService;
        this.clientService = clientService;
    }

    // Prvi endpoint koji pogadja korisnik kada se loguje.
    // Tada zna samo svoje korisnicko ime i lozinku i to prosledjuje na backend.
    @PostMapping("/login")
    public String createAuthenticationToken(@RequestBody UserCredentials authenticationRequest, HttpServletRequest request) {
        if(isUserBlocked(authenticationRequest.getUsername())){
            log.warn("Ip: {}, username: {}, Your account is currently blocked. Try next day again.",request.getRemoteAddr(), authenticationRequest.getUsername());
            return "Your account is currently blocked. Try next day again.";}
        String salt = findSaltForUsername(authenticationRequest.getUsername());
        Authentication authentication = null;
        try {
//             Ukoliko kredencijali nisu ispravni, logovanje nece biti uspesno, desice se AuthenticationException
            authentication = authenticationManager.authenticate(new UsernamePasswordAuthenticationToken(
                    authenticationRequest.getUsername(), authenticationRequest.getPassword().concat(salt)));
            // Ukoliko je autentifikacija uspesna, ubaci korisnika u trenutni security kontekst
            log.info("Ip: {}, username: {}, Login was successful!", request.getRemoteAddr(), authenticationRequest.getUsername());
            SecurityContextHolder.getContext().setAuthentication(authentication);
            refreshMissedPasswordCounter(authenticationRequest.getUsername());
            if(clientService.isPinOk(authenticationRequest.getUsername(), authenticationRequest.getPin()) ||
                    companyService.isPinOk(authenticationRequest.getUsername(), authenticationRequest.getPin()) ||
                    adminService.isPinOk(authenticationRequest.getUsername(), authenticationRequest.getPin())) {
                SecurityContextHolder.getContext().setAuthentication(authentication);
                refreshMissedPasswordCounter(authenticationRequest.getUsername());
            }
            else {
                increaseMissedPasswordCounter(authenticationRequest.getUsername());
                log.warn("Ip: {}, username: {}, Invalid pin.",request.getRemoteAddr(), authenticationRequest.getUsername());
                return "Invalid pin.";
            }
        } catch (AuthenticationException e) {
            increaseMissedPasswordCounter(authenticationRequest.getUsername());
            log.warn("Ip: {}, username: {}, Invalid username, password.",request.getRemoteAddr(), authenticationRequest.getUsername());
            return "Invalid username, password.";
        }

        // Kreiraj token za tog korisnika
        String jwt;
        try {
            CompanyOwner companyOwner;
            if(authentication == null)
                companyOwner = companyService.findByUsername(authenticationRequest.getUsername());
            else
                companyOwner = (CompanyOwner) authentication.getPrincipal();
            if (companyOwner.getForgotten() == 1) {
                companyOwner.setForgotten(2);
                companyService.saveOwner(companyOwner);
            }
            else if (companyOwner.getForgotten() == 2){
                log.warn("Ip: {}, username: {}, Did not changed password first time. If you want to log in, refresh again your password.", request.getRemoteAddr(), authenticationRequest.getUsername());
                return "You did not changed password first time. If you want to log in, refresh again your password.";
            }
            jwt = tokenUtils.generateToken(companyOwner.getUsername(), companyOwner.getRoles());
        } catch (Exception e) {
            try {
                Client client;
                if(authentication == null)
                    client = clientService.findByUsername(authenticationRequest.getUsername());
                else
                    client = (Client) authentication.getPrincipal();
                if (client.getForgotten() == 1) {
                    client.setForgotten(2);
                    clientService.save(client);
                }
                else if (client.getForgotten() == 2){
                    log.warn("Ip: {}, username: {}, Did not changed password first time. If you want to log in, refresh again your password.", request.getRemoteAddr(), authenticationRequest.getUsername());
                    return "You did not changed password first time. If you want to log in, refresh again your password.";
                }
                jwt = tokenUtils.generateToken(client.getUsername(), client.getRoles());
            } catch (Exception e1) {
                Admin admin;
                if(authentication == null)
                    admin = adminService.findByUsername(authenticationRequest.getUsername());
                else
                    admin = (Admin) authentication.getPrincipal();
                jwt = tokenUtils.generateToken(admin.getUsername(), admin.getRoles());
            }
        }

        // Vrati token kao odgovor na uspesnu autentifikaciju
        log.debug("Ip: {}, username: {}, Token successfully generated, JWT: {}", request.getRemoteAddr(), authenticationRequest.getUsername(), jwt);
        return jwt;
    }

    private boolean isUserBlocked(String username) {
        Client client = clientService.findByUsername(username);
        CompanyOwner owner = companyService.findByUsername(username);
        if(client != null && client.getBlockedDate() != null) {
            Calendar c = Calendar.getInstance();
            c.setTime(client.getBlockedDate());
            c.add(Calendar.DATE, 1);
            if (client.isBlocked() && c.getTime().after(new Date()))
                return true;
            else if (client.isBlocked() && c.getTime().before(new Date())) {
                client.setBlocked(false);
                client.setMissedPasswordCounter(0);
                clientService.save(client);
                return false;
            }
        }
        else if(owner != null && owner.getBlockedDate() != null) {
            Calendar c1 = Calendar.getInstance();
            c1.setTime(owner.getBlockedDate());
            c1.add(Calendar.DATE, 1);
            if (owner.isBlocked() && c1.getTime().after(new Date()))
                return true;
            else if (owner.isBlocked() && c1.getTime().before(new Date())) {
                owner.setBlocked(false);
                owner.setMissedPasswordCounter(0);
                companyService.saveOwner(owner);
                return false;
            }
        }

        return false;
    }

    private void increaseMissedPasswordCounter(String username) {
        Client client = clientService.findByUsername(username);
        CompanyOwner owner = companyService.findByUsername(username);
        if(client != null) {
            client.setMissedPasswordCounter(client.getMissedPasswordCounter() + 1);
            if (client.getMissedPasswordCounter() > 5) {
                client.setBlocked(true);
                client.setBlockedDate(new Date());
            }
            clientService.save(client);
        }
        else if(owner != null) {
            owner.setMissedPasswordCounter(owner.getMissedPasswordCounter() + 1);
            if (owner.getMissedPasswordCounter() > 5) {
                owner.setBlocked(true);
                owner.setBlockedDate(new Date());
            }
            companyService.saveOwner(owner);
        }
    }

    private void refreshMissedPasswordCounter(String username) {
        Client client = clientService.findByUsername(username);
        CompanyOwner owner = companyService.findByUsername(username);
        if(client != null) {
            client.setMissedPasswordCounter(0);
            client.setPin("");
            clientService.save(client);
        }
        else if(owner != null) {
            owner.setMissedPasswordCounter(0);
            owner.setPin("");
            companyService.saveOwner(owner);
        }
    }

    private String findSaltForUsername(String username) {
        if(adminService.findByUsername(username) != null)
            return adminService.findByUsername(username).getSalt();
        else if(clientService.findByUsername(username) != null)
            return clientService.findByUsername(username).getSalt();
        else if(companyService.findByUsername(username) != null)
            return companyService.findByUsername(username).getSalt();
        else
            return "";
    }

    @GetMapping(path = "/getAllUsernames")
    public Set<String> getAllUsername() {
        Set<String> usernameList = new HashSet<>();
        usernameList.addAll(adminService.findAllUsernames());
        usernameList.addAll(companyService.findAllUsernames());
        usernameList.addAll(clientService.findAllUsernames());

        return usernameList;
    }

    @GetMapping(path = "/password/blackList/{pass}")
    public ResponseEntity<?> checkPasswordBlackList(@PathVariable String pass) throws URISyntaxException, IOException {
        return clientService.isPasswordInBlackList(pass);
    }

    @GetMapping(path = "/sso/{username}")
    public ResponseEntity<?> sendPasswordlessToken(@PathVariable String username) {
        Set<Role> roles = new HashSet<>();
        Admin admin = adminService.findByUsername(username);
        Client client = clientService.findByUsername(username);
        CompanyOwner owner = companyService.findByUsername(username);
        if(admin != null)
            return emailService.sendPasswordless(admin.getEmail(), tokenUtils.generateToken(username, roles));
        else if(client != null)
            return emailService.sendPasswordless(client.getEmail(), tokenUtils.generateToken(username, roles));
        else if(owner != null)
            return emailService.sendPasswordless(owner.getEmail(), tokenUtils.generateToken(username, roles));

        return new ResponseEntity<>("User with that username does not exist.", HttpStatus.BAD_REQUEST);
    }

    @PostMapping(path = "/login/passwordless")
    public String loginPaswordless(@RequestParam("token") String token) {
        String username = tokenUtils.getUsernameFromToken(token);

        Authentication authentication = new UsernamePasswordAuthenticationToken(username, null);
        SecurityContextHolder.getContext().setAuthentication(authentication);

        String jwt = "";
        Admin admin = adminService.findByUsername(username);
        CompanyOwner owner = companyService.findByUsername(username);
        Client client = clientService.findByUsername(username);
        if(admin != null)
            jwt = tokenUtils.generateToken(admin.getUsername(), admin.getRoles());
        if(client != null)
            jwt = tokenUtils.generateToken(client.getUsername(), client.getRoles());
        if(owner != null)
            jwt = tokenUtils.generateToken(owner.getUsername(), owner.getRoles());

        return jwt;
    }

    @PostMapping(path = "/2factorAuth/pin/send")
    public ResponseEntity<?> sendPinFor2Auth(@RequestBody UserCredentials authenticationRequest) {
        Admin admin = adminService.findByUsername(authenticationRequest.getUsername());
        Client client = clientService.findByUsername(authenticationRequest.getUsername());
        CompanyOwner owner = companyService.findByUsername(authenticationRequest.getUsername());
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        if (admin != null) {
            if (passwordEncoder.matches(authenticationRequest.getPassword().concat(admin.getSalt()), admin.getPassword())) {
                adminService.send2factorAuthPin(admin);
                return new ResponseEntity<>(HttpStatus.OK);
            } else
                increaseMissedPasswordCounter(authenticationRequest.getUsername());
        }
        else if (client != null){
            if (passwordEncoder.matches(authenticationRequest.getPassword().concat(client.getSalt()), client.getPassword())) {
                clientService.send2factorAuthPin(client);
                return new ResponseEntity<>(HttpStatus.OK);
            } else
                increaseMissedPasswordCounter(authenticationRequest.getUsername());
        }
        else if(owner != null) {
            if (passwordEncoder.matches(authenticationRequest.getPassword().concat(owner.getSalt()), owner.getPassword())) {
                companyService.send2factorAuthPin(owner);
                return new ResponseEntity<>(HttpStatus.OK);
            } else
                increaseMissedPasswordCounter(authenticationRequest.getUsername());
        }
        return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
    }

}
