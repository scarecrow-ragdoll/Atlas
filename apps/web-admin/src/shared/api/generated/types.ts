export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  DateTime: { input: string; output: string; }
  UUID: { input: string; output: string; }
};

export type AdminUser = {
  __typename?: 'AdminUser';
  createdAt: Scalars['DateTime']['output'];
  email: Scalars['String']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  role: Scalars['String']['output'];
  updatedAt: Scalars['DateTime']['output'];
};

export type AuthError = {
  __typename?: 'AuthError';
  message: Scalars['String']['output'];
};

export type CreateAdminInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type CreateAdminResult = AuthError | CreateAdminSuccess | ValidationError;

export type CreateAdminSuccess = {
  __typename?: 'CreateAdminSuccess';
  admin: AdminUser;
};

export type CreateUserInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type CreateUserResult = AuthError | CreateUserSuccess | ValidationError;

export type CreateUserSuccess = {
  __typename?: 'CreateUserSuccess';
  user: User;
};

export type DeleteUserResult = AuthError | DeleteUserSuccess;

export type DeleteUserSuccess = {
  __typename?: 'DeleteUserSuccess';
  ok: Scalars['Boolean']['output'];
};

export type LoginAdminInput = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type LoginAdminResult = AuthError | LoginAdminSuccess | ValidationError;

export type LoginAdminSuccess = {
  __typename?: 'LoginAdminSuccess';
  admin: AdminUser;
};

export type LogoutAdminResult = LogoutAdminSuccess;

export type LogoutAdminSuccess = {
  __typename?: 'LogoutAdminSuccess';
  ok: Scalars['Boolean']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createAdmin: CreateAdminResult;
  createUser: CreateUserResult;
  deleteUser: DeleteUserResult;
  loginAdmin: LoginAdminResult;
  logoutAdmin: LogoutAdminResult;
  updateUser: UpdateUserResult;
};


export type MutationCreateAdminArgs = {
  input: CreateAdminInput;
};


export type MutationCreateUserArgs = {
  input: CreateUserInput;
};


export type MutationDeleteUserArgs = {
  id: Scalars['UUID']['input'];
};


export type MutationLoginAdminArgs = {
  input: LoginAdminInput;
};


export type MutationUpdateUserArgs = {
  id: Scalars['UUID']['input'];
  input: UpdateUserInput;
};

export type NotFoundError = {
  __typename?: 'NotFoundError';
  entityType: Scalars['String']['output'];
  id: Scalars['String']['output'];
  message: Scalars['String']['output'];
};

export type PageInfo = {
  __typename?: 'PageInfo';
  endCursor?: Maybe<Scalars['String']['output']>;
  hasNextPage: Scalars['Boolean']['output'];
  hasPreviousPage: Scalars['Boolean']['output'];
  startCursor?: Maybe<Scalars['String']['output']>;
};

export type PaginationInput = {
  after?: InputMaybe<Scalars['String']['input']>;
  before?: InputMaybe<Scalars['String']['input']>;
  first?: InputMaybe<Scalars['Int']['input']>;
  last?: InputMaybe<Scalars['Int']['input']>;
};

export type Query = {
  __typename?: 'Query';
  me?: Maybe<AdminUser>;
  user?: Maybe<User>;
  users: UserConnection;
};


export type QueryUserArgs = {
  id: Scalars['UUID']['input'];
};


export type QueryUsersArgs = {
  pagination?: InputMaybe<PaginationInput>;
};

export type UpdateUserInput = {
  email?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type UpdateUserResult = AuthError | NotFoundError | UpdateUserSuccess | ValidationError;

export type UpdateUserSuccess = {
  __typename?: 'UpdateUserSuccess';
  user: User;
};

export type User = {
  __typename?: 'User';
  createdAt: Scalars['DateTime']['output'];
  email: Scalars['String']['output'];
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  updatedAt: Scalars['DateTime']['output'];
};

export type UserConnection = {
  __typename?: 'UserConnection';
  edges: Array<UserEdge>;
  pageInfo: PageInfo;
  totalCount: Scalars['Int']['output'];
};

export type UserEdge = {
  __typename?: 'UserEdge';
  cursor: Scalars['String']['output'];
  node: User;
};

export type ValidationError = {
  __typename?: 'ValidationError';
  field: Scalars['String']['output'];
  message: Scalars['String']['output'];
};

export type CreateAdminMutationVariables = Exact<{
  input: CreateAdminInput;
}>;


export type CreateAdminMutation = { __typename?: 'Mutation', createAdmin: { __typename: 'AuthError', message: string } | { __typename: 'CreateAdminSuccess', admin: { __typename?: 'AdminUser', id: string, email: string, name: string, role: string, createdAt: string, updatedAt: string } } | { __typename: 'ValidationError', field: string, message: string } };

export type CurrentAdminQueryVariables = Exact<{ [key: string]: never; }>;


export type CurrentAdminQuery = { __typename?: 'Query', me?: { __typename?: 'AdminUser', id: string, email: string, name: string, role: string, createdAt: string, updatedAt: string } | null };

export type LoginAdminMutationVariables = Exact<{
  input: LoginAdminInput;
}>;


export type LoginAdminMutation = { __typename?: 'Mutation', loginAdmin: { __typename: 'AuthError', message: string } | { __typename: 'LoginAdminSuccess', admin: { __typename?: 'AdminUser', id: string, email: string, name: string, role: string, createdAt: string, updatedAt: string } } | { __typename: 'ValidationError', field: string, message: string } };

export type LogoutAdminMutationVariables = Exact<{ [key: string]: never; }>;


export type LogoutAdminMutation = { __typename?: 'Mutation', logoutAdmin: { __typename: 'LogoutAdminSuccess', ok: boolean } };

export type CreateUserMutationVariables = Exact<{
  input: CreateUserInput;
}>;


export type CreateUserMutation = { __typename?: 'Mutation', createUser: { __typename?: 'AuthError', message: string } | { __typename?: 'CreateUserSuccess', user: { __typename?: 'User', id: string, email: string, name: string } } | { __typename?: 'ValidationError', field: string, message: string } };

export type GetUserQueryVariables = Exact<{
  id: Scalars['UUID']['input'];
}>;


export type GetUserQuery = { __typename?: 'Query', user?: { __typename?: 'User', id: string, email: string, name: string, createdAt: string, updatedAt: string } | null };

export type GetUsersQueryVariables = Exact<{
  first?: InputMaybe<Scalars['Int']['input']>;
  after?: InputMaybe<Scalars['String']['input']>;
}>;


export type GetUsersQuery = { __typename?: 'Query', users: { __typename?: 'UserConnection', totalCount: number, edges: Array<{ __typename?: 'UserEdge', cursor: string, node: { __typename?: 'User', id: string, email: string, name: string, createdAt: string } }>, pageInfo: { __typename?: 'PageInfo', hasNextPage: boolean, endCursor?: string | null } } };
