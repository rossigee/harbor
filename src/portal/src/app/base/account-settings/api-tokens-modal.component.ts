import { Component, OnInit, ViewChild } from '@angular/core';
import { MessageHandlerService } from '../../shared/services/message-handler.service';
import { SessionService } from '../../shared/services/session.service';
import { UserService } from 'ng-swagger-gen/services/user.service';
import { InlineAlertComponent } from '../../shared/components/inline-alert/inline-alert.component';

@Component({
    selector: 'api-tokens-modal',
    templateUrl: 'api-tokens-modal.component.html',
    styleUrls: ['./api-tokens-modal.component.scss', '../../common.scss'],
    standalone: false,
})
export class ApiTokensModalComponent implements OnInit {
    opened = false;
    staticBackdrop = true;

    tokens: any[] = [];
    selectedTokens: any[] = [];
    tokenLoading = false;
    showCreateTokenModal = false;
    createdTokenSecret: string;
    newTokenForm: any = { name: '', description: '' };

    @ViewChild(InlineAlertComponent)
    inlineAlertComponent: InlineAlertComponent;

    constructor(
        private session: SessionService,
        private msgHandler: MessageHandlerService,
        private userService: UserService
    ) {}

    ngOnInit(): void {
        this.loadTokens();
    }

    open() {
        this.opened = true;
        this.loadTokens();
    }

    close() {
        this.opened = false;
    }

    loadTokens(): void {
        this.tokenLoading = true;
        // TODO: Implement API tokens loading from API
        // this.userService.listUserTokens().subscribe(
        //     (data) => {
        //         this.tokens = data;
        //         this.tokenLoading = false;
        //     },
        //     (error) => {
        //         this.msgHandler.handleError(error);
        //         this.tokenLoading = false;
        //     }
        // );
    }

    openCreateTokenModal(): void {
        this.showCreateTokenModal = true;
        this.newTokenForm = { name: '', description: '' };
    }

    closeCreateTokenModal(): void {
        this.showCreateTokenModal = false;
        this.createdTokenSecret = null;
    }

    createToken(): void {
        if (!this.newTokenForm.name) {
            this.msgHandler.showError('Token name is required', {});
            return;
        }

        this.tokenLoading = true;
        // TODO: Implement token creation API call
        // this.userService.createToken(this.newTokenForm).subscribe(
        //     (response) => {
        //         this.createdTokenSecret = response.secret;
        //         this.tokens.push(response);
        //         this.tokenLoading = false;
        //     },
        //     (error) => {
        //         this.msgHandler.handleError(error);
        //         this.tokenLoading = false;
        //     }
        // );
    }

    copyTokenSecret(): void {
        const copyInput = document.createElement('textarea');
        copyInput.value = this.createdTokenSecret;
        document.body.appendChild(copyInput);
        copyInput.select();
        document.execCommand('copy');
        document.body.removeChild(copyInput);
        this.msgHandler.showSuccess('Token copied to clipboard');
    }

    revokeToken(tokenId: any): void {
        // TODO: Implement token revocation
    }

    deleteToken(tokenId: any): void {
        this.tokens = this.tokens.filter(t => t.id !== tokenId);
    }

    formatScope(scope: any): string {
        if (typeof scope === 'string') {
            return scope;
        }
        if (Array.isArray(scope)) {
            return scope.join(', ');
        }
        return JSON.stringify(scope);
    }
}
