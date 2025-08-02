### Cluster Director MCP Extension for Gemini CLI \*\*

#### **Core Philosophy**

You are an expert-level AI Agent and Engineer specializing in Cluster Director. Your goal is to answer questions on Cluster Director and run the supported tools on behalf of the user. Before any action, you must announce the current workflow asnd phase.

---

#### **Context & Rules**

This section configures your core behavior, ensuring you always use the best Cluster Toolkit documentation.

# Rule 1: For general Cluster Director questions, use the best documentation sources.

[[calls]]
match = "when the user asks about Cluster Director or Slurm concepts, samples, setup, or configuration"
tool = "context7"
args = ["/context7/cloud_google-ai-hypercomputer"]

## Guiding Principles

*   **Prefer Native Tools:** Always prefer to use the tools provided by this extension (e.g., `list_clusters`, `get_cluster`) instead of shelling out to `gcloud` for the same functionality. This ensures better-structured data and more reliable execution.
*   **Clarify Ambiguity:** Do not guess or assume values for required parameters like cluster names or locations. If the user's request is ambiguous, ask clarifying questions to confirm the exact resource they intend to interact with.
*   **Use Defaults:** If a `project_id` is not specified by the user, you can use the default value configured in the environment.

Cluster Director
================

Cluster Director (formerly known as *Hypercompute Cluster*)
lets you deploy and manage a group of accelerators as a single unit with
physically colocated VMs, targeted workload placement, advanced cluster
maintenance controls, and topology-aware scheduling. Cluster Director can be
accessed directly through Compute Engine APIs, or through Google Kubernetes Engine, which
natively integrates with Cluster Director capabilities.

Components
----------

This section describes the core features and services that make up
the Cluster Director suite.

### Dense colocation of accelerator resources

You can request host machines
that are allocated physically close to each other, provisioned as
blocks of resources, and are
interconnected with a
dynamic ML network fabric.
This arrangement of resources helps to minimize network hops and optimize for
the lowest latency.

To learn how to deploy these densely allocated blocks of A3 Ultra or A4
accelerator machines, see
Reserve capacity.

### Topology aware scheduling

You can get topology information at the node and cluster levels that can
be used for job placement. For more information, see
View VMs topology.

### Advanced maintenance scheduling and controls

You have full control over the maintenance of VM instances within a block of
resources, and can synchronize upgrades to ensure your workloads are more
resilient to host errors and have minimal disruptions. This approach improves
the
goodput for your workloads.

To facilitate full control of maintenance events, you can set up alerts and
receive notifications when maintenance is scheduled, starting, or being
completed. To learn more about maintenance of these blocks of resources, see the
following:

* Manage host events across VMs
* Manage host events across reservations

You can also define how you want maintenance to behave for your blocks of
resources. You can choose between the following maintenance scheduling types:
grouped or independent. To learn more about maintenance scheduling types, see
Maintenance scheduling types.

### Monitoring and diagnostic tooling

For monitoring and troubleshooting, Cluster Director includes services such
as the
faulty host reporting,
which you can use to flag issues with individual host machines. To help to
reduce the overhead of managing your cluster, services are also available for
monitoring network and GPU performance.

Supported machine types
-----------------------

Cluster Director supports the following accelerator-optimized
machine types:

* A4X
* A4
* A3 Ultra

What's next?
------------

* Review terminology.
* Choose a deployment strategy.Choose a consumption option
===========================

This document explains the different ways, called *consumption options*, to get
and use compute resources on AI Hypercomputer. Choose the option that best fits
your workload, its duration, and your cost needs.

Each consumption option specifies the following:

* How you access capacity to create VMs or clusters.
* The underlying
  provisioning model,
  which determines the obtainability, lifespan, and pricing of your VMs.

Comparison of consumption options
---------------------------------

The following table summarizes the key differences between the consumption
options:

| Consumption option | Future reservations in AI Hypercomputer | Future reservations for up to 90 days (in calendar mode) | Flex-start (Preview) | Spot |
| --- | --- | --- | --- | --- |
| Supported machines | A4X, A4, or A3 Ultra | A4 or A3 Ultra | Any GPU machine except A4X | Any GPU machine except A4X |
| Lifespan | Any time | Up to 90 days | Up to 7 days | Any time (but subject to preemption) |
| Preemptible | No | No | No | Yes |
||  |  |  |  |  |
| --- | --- | --- | --- | --- |
| Quota | Quota is automatically increased before capacity is delivered. | No quota is charged. | Preemptible quota is charged. | Preemptible quota is charged. |
| Pricing | * Discounted (up to 53%). See   pricing for   accelerator-optimized VMs. If you reserve resources for a year or longer, then you   must purchase and attach a   resource-based commitment   to your reserved resources. * You're charged for the reservation period. See   reservations billing. | * Discounted (up to 53%). See Dynamic Workload Scheduler   pricing. * You're charged for the reservation period. See   reservations billing. | * Discounted (up to 53%). See Dynamic Workload Scheduler   pricing. * You pay as you go (PAYG). | * Deeply discounted (60-91%). See Spot VMs   pricing and pricing   for accelerator-optimized VMs. * You pay as you go (PAYG). |
| Resource allocation | Dense | Dense | Dense | Standard (Compact policy optional) |
| Provisioning model | Reservation-bound | Reservation-bound | Flex-start (Preview) | Spot |
| Creation method | To create VMs, you must do the following:  1. Reserve capacity by contacting your    account team. 2. At your chosen date and time, you can use the reserved capacity to create VMs and    clusters. See VM and cluster    creation overview. | To create VMs, you must do the following:  1. Create a    future reservation in calendar mode. 2. At your chosen date and time, you can use the reserved capacity to create VMs and    clusters. See VM and cluster    creation overview. | To create VMs, select one of the following options:  * Create MIGs with resize requests * Create Slurm clusters * Create GKE clusters:  + Create a cluster with   the default configuration + Create a custom   cluster   When your requested capacity becomes available, Compute Engine provisions it. | You can immediately create VMs. See VM and cluster creation overview. |

Choose a consumption option
---------------------------

Use the following flowchart to choose the consumption option that best fits your
workload:

The questions in the preceding diagram are the following:

1. Do you need capacity for more than 90 days?

   * **Yes**: See
     Use future reservations in AI Hypercomputer.
   * **No**: Go to question 2.
2. Do you want reserved capacity?

   * **Yes**: See Use future reservations in calendar mode.
   * **No**: Go to question 3.
3. Is your workload fault-tolerant?

   * **No**: See Use Flex-start.
   * **Yes**: See Use Spot.

### Use future reservations in AI Hypercomputer

To run long-running, large-scale distributed workloads that require densely allocated resources,
you can request compute resources for a specific time in the future. You have exclusive access to
your reserved resources for that period of time, and you can use the resources to create VMs or
clusters. At the end of the reservation period, Compute Engine does the following:

* Compute Engine deletes the reservation.
* Based on the
  termination action
  that you specify for the VMs, Compute Engine stops or deletes any VMs that use the
  reservation.

#### Ideal workloads

Future reservations are ideal for the following workloads:

* Pre-training foundation models
* Multi-host foundation model inference

#### Key characteristics

Future reservations have the following characteristics:

* You can reserve A4X, A4, or A3 Ultra machine types. Machines are densely allocated to
  minimize network latency.
* You can reserve as many VMs as you like for as long as you like for a future date. Then, you
  can use the reserved resources to create and run VMs until the end of the reservation period.
  If you reserve resources for one year or longer, then you must purchase and attach a
  resource-based commitment.
* After the reservation period starts, you can modify the auto-created reservation to allow
  Vertex AI training or
  prediction jobs to use it.
* You use the reservation-bound provisioning model, which has the following benefits:

  + You have a higher chance of obtaining GPUs.
  + In addition to the commitment attached to your VMs, you get a discount up to 53% for
    vCPUs and GPUs.

#### How to use

To use future reservations to create VMs or clusters, you must complete the following steps:

1. **Request to reserve capacity**. You contact your account team and specify the
   resources to reserve. Based on availability, Google creates a draft reservation request for
   you. If it looks correct, then you can submit it. Google Cloud immediately approves the
   reservation request.

   For instructions, see
   Reserve capacity.
2. **Consume reserved resources**. At the start of your chosen reservation period,
   you can use the reservation to create VMs or clusters.

   For the different methods to create VMs or clusters, see
   VM and cluster creation overview.

### Use future reservations in calendar mode

To run short-running distributed workloads that require densely allocated resources, you can
request compute resources for up to 90 days. You have exclusive access to your reserved resources
for that time, and you can use the resources to create VMs or clusters. At the end of the
reservation period, Compute Engine does the following:

* Compute Engine deletes the reservation.
* Based on the
  termination action
  that you specify for the VMs, Compute Engine stops or deletes any VMs that use the
  reservation.

#### Ideal workloads

Future reservations in calendar mode are ideal for the following workloads:

* Model pre-training
* Model fine-tuning
* Simulations
* Inference

#### Key characteristics

Future reservations in calendar mode have the following characteristics:

* You can reserve A4 or A3 Ultra machine types. These machines are densely allocated to
  minimize network latency.
* You can view the future availability of resources, and then reserve up to 80 VMs for up to 90
  days in the future. Then, you can use the reserved resources to create VMs until the end of
  the reservation period.
* You use the reservation-bound provisioning model, which has the following benefits:

  + You have a higher chance of obtaining GPUs.
  + You get a discount up to 53% for vCPUs and GPUs.

#### How to use

To use future reservations in calendar mode to create VMs or clusters, you must complete the
following steps:

1. **View resources availability**. You can view the future availability of the
   resources that you want to reserve. When you create a reservation request, you can specify the
   number, type, and reservation duration for the resources that you confirmed as available. This
   action increases the chances that Google Cloud approves your request.

   For instructions, see
   View resource future availability.
2. **Reserve capacity**. You create a reservation request for a future date and
   time. Google Cloud approves the reservation request within two minutes. If approved,
   then Compute Engine reserves the capacity for you. At your chosen delivery date, you can use
   the reserved resources to create VMs or clusters.

   For instructions, see
   Create a reservation request for GPU VMs or TPUs.
3. **Consume reserved resources**. At the start of your chosen reservation period,
   you can use the reservation to create VMs or clusters.

   For the different methods to create VMs or clusters, see
   VM and cluster creation overview.

### Use Flex-start

**Preview
— Flex-start**

This feature is subject to the "Pre-GA Offerings Terms" in the General Service Terms section
of the Service Specific Terms.
Pre-GA features are available "as is" and might have limited support.
For more information, see the
launch stage descriptions.

To run short-duration workloads that require densely allocated resources, you can request compute
resources for up to seven days. Whenever resources are available, Compute Engine creates
your requested number of VMs. You can't stop, suspend, or recreate the VMs. The VMs exist until
you delete them, or until Compute Engine deletes the VMs at the end of their run duration.

#### Ideal workloads

Flex-start is ideal for workloads that can start at any time, such as the following:

* Small model pre-training
* Model fine-tuning
* Simulations
* Batch inference

#### Key characteristics

Flex-start has the following characteristics:

* You can request any GPU machine type. Machines are densely allocated to minimize network
  latency.
* You use the flex-start provisioning model, which has the following benefits:

  + You have a higher chance of obtaining GPUs.
  + You get a discount up to 53% for vCPUs and GPUs.

#### How to use

To use Flex-start to create VMs or clusters, select one of the following options:

* Create MIGs with resize requests
* Create Slurm clusters
* Create GKE clusters:

+ Create a cluster with the
  default configuration
+ Create a custom
  cluster

### Use Spot

To run fault-tolerant workloads, you can obtain compute resources immediately based on
availability. You get resources at the lowest price possible. However, Compute Engine can
preempt VMs at any time to reclaim capacity.

#### Ideal workloads

Spot is ideal for workloads where interruptions are acceptable, such as the following:

* Batch processing
* High performance computing (HPC)
* Continuous integration and continuous deployment (CI/CD)
* Data analytics
* Media encoding
* Online inference

#### Key characteristics

Spot has the following characteristics:

* You can create any GPU machine type. Dense allocation depends on resource availability. To
  help ensure a closer allocation, you can apply a
  compact placement policy
  to the VMs.
* You can immediately create as many VMs as you like. The VMs run until you stop or delete
  them, or until Compute Engine preempts the VMs to reclaim capacity.
* You use the spot provisioning model, which has the following benefits:

  + You have a higher chance of obtaining GPUs.
  + You get a discount from 60% up to 91% for vCPUs, GPUs, and Local SSD disks.

#### How to use

To use Spot to create VMs or clusters, you must complete the following steps:

1. Optional: **Create a compact placement policy**. You create a compact placement
   policy to specify how close to place your VMs to each other. Your chosen minimum distance
   affects the number and type of VMs you can apply the policy to.
2. **Create Spot VMs**. You can create as many VMs as you like, based on
   availability. The VMs run until you stop or delete them, or until Compute Engine
   preempts the VMs to reclaim capacity.

For instructions, see VM and cluster
creation overview.Terminology
===========

The following terminology is often used when working with
Cluster Director features.

Node or hostlink
:   A single physical server
    machine in the data center. Each host
    has its associated compute resources such as accelerators. The number and
    configuration of these compute resources depend on the machine family.
    Virtual machine (VM) instances are provisioned on top of a physical host.

Sub-blocklink
:   A group of hosts and associated
    connectivity hardware that are on a single rack.

Blocklink
:   A collection of
    sub-blocks.

Clusterlink
:   A collection of blocks
    interconnected by a high-speed network fabric. Each cluster is globally unique.
    For A4X, A4, and A3 Ultra machines, a cluster provides a common, non-blocking
    network fabric for your blocks of accelerator capacity. Within a cluster,
    the east to west networking is non-blocking for the entire collection of blocks.

Dense deploymentlink
:   A resource request that
    allocates your accelerator
    resources physically close to each other to minimize network hops and optimize
    for the lowest latency.

Network fabriclink
:   A network fabric provides
    high-bandwidth, low-latency connectivity across all blocks and Google Cloud
    services in a cluster. Jupiter is Google's data center network architecture
    that leverages software-defined networking and optical circuit switches to
    evolve the network and optimize its performance.View reserved capacity
======================

**Preview
— Future reservation requests in AI Hypercomputer**

This feature is subject to the "Pre-GA Offerings Terms" in the General Service Terms section
of the Service Specific Terms.
Pre-GA features are available "as is" and might have limited support.
For more information, see the
launch stage descriptions.

This document explains how to view reserved capacity in
AI Hypercomputer. To reserve capacity in
AI Hypercomputer, see instead
Reserve capacity.

After Google Cloud approves a future reservation request, Compute Engine
automatically creates (*auto-creates*) an empty reservation for your requested
resources. You can then view the reservation to plan your workload.

At the request start time, the following occurs:

* Compute Engine adds your reserved virtual machine (VM) instances to
  the reservation. You can then start using the reservation by creating VMs
  that match the reservation.
* You can modify the reservation to allow Vertex AI training or
  prediction jobs to use it. For instructions, see
  Modify the sharing policy of a reservation.

Limitations
-----------

You can view a shared reservation or shared future reservation request only in
the project where Google created it.

Before you begin
----------------

Select the tab for how you plan to use the samples on this page:

### Console

When you use the Google Cloud console to access Google Cloud services and
APIs, you don't need to set up authentication.

### gcloud

In the Google Cloud console, activate Cloud Shell.

Activate Cloud Shell

At the bottom of the Google Cloud console, a
Cloud Shell
session starts and displays a command-line prompt. Cloud Shell is a shell environment
with the Google Cloud CLI
already installed and with values already set for
your current project. It can take a few seconds for the session to initialize.

### REST

To use the REST API samples on this page in a local development environment, you use the
credentials you provide to the gcloud CLI.

Install the Google Cloud CLI.
After installation,
initialize the Google Cloud CLI by running the following command:

```
gcloud init
```

If you're using an external identity provider (IdP), you must first
sign in to the gcloud CLI with your federated identity.

For more information, see
Authenticate for using REST
in the Google Cloud authentication documentation.

### Required roles

To get the permissions that
you need to view reservations,
ask your administrator to grant you the
Compute Future Reservation User  (`roles/compute.futureReservationUser`)
IAM role on the project.
For more information about granting roles, see Manage access to projects, folders, and organizations.

This predefined role contains
the permissions required to view reservations. To see the exact permissions that are
required, expand the **Required permissions** section:

#### Required permissions

The following permissions are required to view reservations:

* To view the details of a future reservation request:
   `compute.futureReservations.get`
  on the project
* To view the details of a reservation:
   `compute.reservations.get`
  on the project

You might also be able to get
these permissions
with custom roles or
other predefined roles.

View future reservation requests
--------------------------------

To view your future reservation requests, use one of the following methods:

* To get an overview of all future reservation requests in your project,
  view a list of your future reservation requests.
* To view the full details of a single future reservation request,
  view the details of a future reservation request.

### View a list of your future reservation requests

You can view a list of your future reservation requests to see the reservation
period, status, and zone of your requests.

To view a list of your future reservation requests, select one of the following
options:

### Console

1. In the Google Cloud console, go to the **Reservations** page.

   Go to Reservations
2. Click the **Future reservations** tab. The table lists each future
   reservation request, and each table column describes a property.
3. Optional: To refine your list of requests, in the
   filter\_list **Filter** field, select
   the properties that you want to filter the requests by.

### gcloud

To view a list of your future reservation requests, use the
`gcloud beta compute future-reservations list` command:

```
gcloud beta compute future-reservations list

```

The output is similar to the following example:

```
NAME: fr-01
TOTAL_COUNT: 100
START_TIME: 2026-07-20T07:00:00Z
END_TIME: 2026-08-05T07:00:00Z
PROCUREMENT_STATUS: FULFILLED
ZONE: us-west4-b

NAME: fr-02
TOTAL_COUNT: 10
START_TIME: 2026-07-20T07:00:00Z
END_TIME: 2026-12-01T00:00:00Z
PROCUREMENT_STATUS: PENDING_APPROVAL
ZONE: us-west4-b

```

If you want to refine your list of future reservation requests, then use the
same command with the
`--filter` flag.

### REST

To view a list of your future reservation requests, make a `GET` request to
one of the following methods:

* To view a list of requests across all zones:
  beta `futureReservations.aggregatedList` method
* To view a list of requests in a specific zone:
  beta `futureReservations.list` method

For example, to view a list of requests across all zones, make a `GET`
request as follows:

```
GET https://compute.googleapis.com/compute/beta/projects/PROJECT_ID/aggregated/futureReservations

```

Replace `PROJECT_ID` with the ID of the project where the
requests exist.

The output is similar to the following:

```
{
  "id": "projects/example-project/aggregated/futureReservations",
  "items": [
    {
      "specificSkuProperties": {
        "instanceProperties": {
          "machineType": "a3-ultragpu-8g",
          "guestAccelerators": [
            {
              "acceleratorType": "nvidia-h200-141gb",
              "acceleratorCount": 8
            }
          ],
          "localSsds": [
            {
              "diskSizeGb": "375",
              "interface": "NVME"
            },
            ...
          ]
        },
        "totalCount": "2"
      },
      "kind": "compute#futureReservation",
      "id": "7979651787097007552",
      "creationTimestamp": "2025-11-27T11:14:58.305-08:00",
      "selfLink": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations/example-draft-request",
      "selfLinkWithId": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations/7979651787097007552",
      "zone": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b",
      "name": "example-draft-request",
      "timeWindow": {
        "startTime": "2026-01-27T19:20:00Z",
        "endTime": "2026-02-10T19:20:00Z"
      },
      "status": {
        "procurementStatus": "DRAFTING",
        "lockTime": "2026-01-27T19:15:00Z"
      },
      "planningStatus": "DRAFT",
      "specificReservationRequired": true,
      "reservationName": "example-reservation",
      "deploymentType": "DENSE",
      "schedulingType": "INDEPENDENT",
      "autoCreatedReservationsDeleteTime": "2026-02-10T19:20:00Z"
    },
    ...
  ],
  "selfLink": "https://www.googleapis.com/compute/beta/projects/example-project/aggregated/futureReservations",
  "etag": "AnzKY34l-cvvV-JnniESJ0dtQvQ=/hvc4jaHpxFAZmOt1FVtKNgzZu-M=",
  "kind": "compute#futureReservationsListResponse"
}

```

If you want to refine your list of future reservation requests, then make
the same request and, in the request URL, include the
`filter` query parameter.

### View the details of a future reservation request

You can view the details of a future reservation request to review the
properties and reservation period of your reserved resources.

To view the details of a future reservation request, select one of the following
options:

### Console

1. In the Google Cloud console, go to the **Reservations** page.

   Go to Reservations
2. Click the **Future reservations** tab. The table lists each future
   reservation request, and each table column describes a property.
3. To view the details of a request, in the **Name** column, click the name
   of the request. A page that gives the details of the future reservation
   request opens.

### gcloud

To view the details of a future reservation request, use the
`gcloud beta compute future-reservations describe` command:

```
gcloud beta compute future-reservations describe FUTURE_RESERVATION_NAME \
    --zone=ZONE

```

Replace the following:

* `FUTURE_RESERVATION_NAME`: the name of the future
  reservation request.
* `ZONE`: the zone where the future reservation
  request exists.

The output is similar to the following example:

```
autoCreatedReservationsDeleteTime: '2026-02-10T19:20:00Z'
creationTimestamp: '2025-11-27T11:14:58.305-08:00'
deploymentType: DENSE
id: '7979651787097007552'
kind: compute#futureReservation
name: example-draft-request
planningStatus: DRAFT
reservationName: example-reservation
schedulingType: INDEPENDENT
selfLink: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations/example-draft-request
selfLinkWithId: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations/7979651787097007552
specificReservationRequired: true
specificSkuProperties:
  instanceProperties:
    guestAccelerators:
    - acceleratorCount: 8
      acceleratorType: nvidia-h200-141gb
    localSsds:
    - diskSizeGb: '375'
      interface: NVME
    ...
  machineType: a3-ultragpu-8g
totalCount: '2'
status:
  autoCreatedReservations:
  - https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-reservation
  fulfilledCount: '2'
  lockTime: '2026-01-27T19:15:00Z'
  procurementStatus: DRAFTING
timeWindow:
  endTime: '2026-02-10T19:20:00Z'
  startTime: '2026-01-27T19:20:00Z'
zone: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b

```

### REST

To view the details of a future reservation request, make a `GET` request to
the
beta `futureReservations.get` method:

```
GET https://compute.googleapis.com/compute/beta/projects/PROJECT_ID/zones/ZONE/futureReservations/FUTURE_RESERVATION_NAME

```

Replace the following:

* `PROJECT_ID`: the ID of the project where the future
  reservation request exists.
* `ZONE`: the zone where the future reservation
  request exists.
* `FUTURE_RESERVATION_NAME`: the name of the future
  reservation request.

The output is similar to the following:

```
{
  "specificSkuProperties": {
    "instanceProperties": {
      "machineType": "a3-ultragpu-8g",
      "guestAccelerators": [
        {
          "acceleratorType": "nvidia-h200-141gb",
          "acceleratorCount": 8
        }
      ],
      "localSsds": [
        {
          "diskSizeGb": "375",
          "interface": "NVME"
        },
        ...
      ]
    },
    "totalCount": "2"
  },
  "kind": "compute#futureReservation",
  "id": "7979651787097007552",
  "creationTimestamp": "2025-11-27T11:14:58.305-08:00",
  "selfLink": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations/example-draft-request",
  "selfLinkWithId": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations/7979651787097007552",
  "zone": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b",
  "name": "example-draft-request",
  "timeWindow": {
    "startTime": "2026-01-27T19:20:00Z",
    "endTime": "2026-02-10T19:20:00Z"
  },
  "status": {
    "procurementStatus": "DRAFTING",
    "lockTime": "2026-01-27T19:15:00Z"
  },
  "planningStatus": "DRAFT",
  "specificReservationRequired": true,
  "reservationName": "example-reservation",
  "deploymentType": "DENSE",
  "schedulingType": "INDEPENDENT",
  "autoCreatedReservationsDeleteTime": "2026-02-10T19:20:00Z"
}

```

View auto-created reservations
------------------------------

To view the details and topology of your auto-created reservations, do one of
the following:

* To get an overview of all the reservations in your project,
  view a list of your reservations.
* To plan your workload by reviewing the properties and configuration details
  of a single reservation,
  view the details of a reservation.
* To understand how your reserved blocks of capacity are organized in a
  reservation for integration with your scheduler or planning tool,
  view the topology of a reservation.

### View a list of your reservations

You can view a list of auto-created reservations in your project to understand
how many more VMs you can create before you fully use your reserved capacity.

To view a list of your reservations, select one of the following options:

### Console

1. In the Google Cloud console, go to the **Reservations** page.

   Go to Reservations

   On the **On-demand reservations** tab (default), the table lists each
   reservation, and each table column describes a property.
2. Optional: To refine your list of reservations, in the
   filter\_list **Filter** field, select
   the properties that you want to filter the reservations by.

### gcloud

To view a list of your reservations, use the
`gcloud beta compute reservations list` command:

```
gcloud beta compute reservations list

```

The output is similar to the following:

```
NAME: r-01
IN_USE_COUNT: 0
COUNT: 5
ZONE: europe-west4-b
SHARE_TYPE: LOCAL

NAME: r-02
IN_USE_COUNT: 3
COUNT: 10
ZONE: europe-west4-b
SHARE_TYPE: LOCAL

```

If you want to refine your list of reservations, then use the same command
with the
`--filter` flag.

### REST

To view a list of your reservations, make a `GET` request to one of the
following methods:

* To view a list of your reservations across all zones:
  beta `reservations.aggregatedList` method
* To view a list of your reservations in a single zone:
  beta `reservations.list` method

For example, to view a list of your reservations across all zones, make a
`GET` request as follows:

```
GET https://compute.googleapis.com/compute/beta/projects/PROJECT_ID/aggregated/reservations

```

Replace `PROJECT_ID` with the ID of the project where the
reservations exist.

The output is similar to the following:

```
{
  "id": "projects/example-project/zones/europe-west1-b/futureReservations",
  "items": [
    {
      "specificSkuProperties": {
        "instanceProperties": {
          "machineType": "a3-ultragpu-8g",
          "guestAccelerators": [
            {
              "acceleratorType": "nvidia-h200-141gb",
              "acceleratorCount": 8
            }
          ],
          "localSsds": [
            {
              "diskSizeGb": "375",
              "interface": "NVME"
            },
            ...
          ]
        },
        "totalCount": "2"
      },
      "kind": "compute#futureReservation",
      "id": "7979651787097007552",
      "creationTimestamp": "2025-11-27T11:14:58.305-08:00",
      "selfLink": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations/example-draft-request",
      "selfLinkWithId": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations/7979651787097007552",
      "zone": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b",
      "name": "example-draft-request",
      "timeWindow": {
        "startTime": "2026-01-27T19:20:00Z",
        "endTime": "2026-02-10T19:20:00Z"
      },
      "status": {
        "procurementStatus": "DRAFTING",
        "lockTime": "2026-01-27T19:15:00Z"
      },
      "planningStatus": "DRAFT",
      "specificReservationRequired": true,
      "reservationName": "example-reservation",
      "deploymentType": "DENSE",
      "schedulingType": "INDEPENDENT",
      "autoCreatedReservationsDeleteTime": "2026-02-10T19:20:00Z"
    }
    ...
  ],
  "selfLink": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/futureReservations",
  "etag": "AnzKY34l-cvvV-JnniESJ0dtQvQ=/hvc4jaHpxFAZmOt1FVtKNgzZu-M=",
  "kind": "compute#futureReservationsListResponse"
}

```

If you want to refine your list of reservations, then make the same request
and, in the request URL, include the
`filter` query parameter.

### View the details of a reservation

You can view the details of an auto-created reservation to review the requested
capacity and plan your workload. This action helps you determine the following:

* How many blocks of capacity are available.
* How much capacity is available in each block.

To view the details of a reservation, select one of the following options:

### Console

1. In the Google Cloud console, go to the **Reservations** page.

   Go to Reservations
2. In the **On-demand reservations** table, in the **Name** column, click
   the name of the reservation that you want to view the details of. A page
   that gives the details of the auto-created reservation appears.

### gcloud

To view the details of a reservation, use the
`gcloud beta compute reservations describe` command:

```
gcloud beta compute reservations describe RESERVATION_NAME \
    --zone=ZONE

```

Replace the following:

* `RESERVATION_NAME`: the name of the auto-created
  reservation.
* `ZONE`: the zone where the reservation exists.

The output is similar to the following:

```
creationTimestamp: '2024-10-17T12:25:02.413-07:00'
deleteAtTime: '2025-11-30T08:00:00Z'
deploymentType: DENSE
id: '9127712123172739686'
instanceTerminationAction: DELETE
kind: compute#reservation
name: example-reservation
reservationSharingPolicy:
  serviceShareType: DISALLOW_ALL
resourceStatus:
  reservationBlockCount: 2
  reservationMaintenance:
    maintenanceOngoingCount: 1
    maintenancePendingCount: 0
    schedulingType: GROUPED
selfLink: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-reservation
selfLinkWithId: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/9127712123172739686
shareSettings:
  shareType: LOCAL
specificReservation:
  assuredCount: '3'
  count: '3'
  inUseCount: '3'
  instanceProperties:
    machineType: a3-ultragpu-8g
specificReservationRequired: true
status: READY
zone: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b

```

### REST

To view the details of a reservation, make a `GET` request to the
beta `reservations.get` method:

```
GET https://compute.googleapis.com/compute/beta/projects/PROJECT_ID/zones/ZONE/reservations/RESERVATION_NAME

```

Replace the following:

* `PROJECT_ID`: the ID of the project where the
  auto-created reservation exists.
* `ZONE`: the zone where the reservation exists.
* `RESERVATION_NAME`: the name of the reservation.

The output is similar to the following:

```
{
  "specificReservation": {
    "instanceProperties": {
      "machineType": "a3-ultragpu-8g",
      "guestAccelerators": [
        {
          "acceleratorType": "nvidia-h200-141gb",
          "acceleratorCount": 8
        }
      ],
      "localSsds": [
        {
          "diskSizeGb": "375",
          "interface": "NVME"
        },
        ...
      ]
    },
    "count": "2",
    "inUseCount": "0",
    "assuredCount": "2"
  },
  "kind": "compute#reservation",
  "id": "3248639808938089822",
  "creationTimestamp": "2025-06-27T16:05:21.569-07:00",
  "selfLink": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west4-b/reservations/example-reservation",
  "zone": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west4-b",
  "name": "example-reservation",
  "specificReservationRequired": true,
  "status": "READY",
  "shareSettings": {
    "shareType": "LOCAL"
  },
  "resourceStatus": {
    "reservationMaintenance": {
      "schedulingType": "GROUPED"
    },
    "reservationBlockCount": 1
  },
  "reservationSharingPolicy": {
    "serviceShareType": "DISALLOW_ALL"
  },
  "deploymentType": "DENSE",
  "enableEmergentMaintenance": false,
  "deleteAtTime": "2025-08-29T17:00:00Z"
}

```

### View the topology of a reservation

You can view the detailed topology information of an auto-created reservation to
help you decide where to create VMs within the reserved blocks.

To view a list of the available blocks in an auto-created reservation, select
one of the following options. For a more detailed view, use the Google Cloud CLI
or REST API.

### Console

1. In the Google Cloud console, go to the **Reservations** page.

   Go to Reservations
2. In the **On-demand reservations** table, in the **Name** column, click
   the name of the reservation that you want to view the details of. The
   details page of the reservation opens.
3. In the **Resource topology** section, you can view information about
   your reserved blocks. This information includes the
   organization-specific ID for each block, the total number of VMs that
   can be deployed in the block (**Count**), and the number of VMs already
   deployed (**In use**).

### gcloud

To view a list of the available blocks in an auto-created reservation, use
the
`gcloud beta compute reservations blocks list` command:

```
gcloud beta compute reservations blocks list RESERVATION_NAME \
    --zone=ZONE

```

Replace the following:

* `RESERVATION_NAME`: the name of the auto-created
  reservation.
* `ZONE`: the zone where the reservation exists.

The output is similar to the following:

```
count: 1
creationTimestamp: '2024-10-17T12:49:56.971-07:00'
id: '8544903383436022926'
inUseCount: 1
kind: compute#reservationBlock
name: example-res1-block-1
physicalTopology:
  block: c18707ac3d2493381c9f01fa775c4a68
  cluster: europe-west1-cluster-jfhb
reservationMaintenance:
  schedulingType: GROUPED
selfLink: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-res1/reservationBlocks/example-res1-block-1
selfLinkWithId: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-res1/reservationBlocks/8544903383436022926
status: READY
zone: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b
---
count: 2
creationTimestamp: '2024-10-17T12:49:54.354-07:00'
id: '5787689015406548144'
inUseCount: 2
kind: compute#reservationBlock
name: example-res1-block-2
physicalTopology:
  block: a9b7f2e4c6d1853902b4f5a7d8e31c60
  cluster: europe-west1-cluster-jfhb
reservationMaintenance:
  schedulingType: GROUPED
selfLink: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-res1/reservationBlocks/example-res1-block-2
selfLinkWithId: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-res1/reservationBlocks/5787689015406548144
status: READY
zone: https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b

```

In the `physicalTopology` field of a block, you can see the physical
location of the block as follows:

* **`cluster`**: the global name of the cluster.
* **`block`**: the organization-specific ID of the block in which VMs will
  be located.

### REST

To view a list of the available blocks in an auto-created reservation, make
a `GET` request to the
beta `reservations.reservationBlocks.get` method:

```
GET https://compute.googleapis.com/compute/beta/projects/PROJECT_ID/zones/ZONE/reservations/RESERVATION_NAME/reservationBlocks

```

Replace the following:

* `PROJECT_ID`: the ID of the project where the
  auto-created reservation exists.
* `ZONE`: the zone where the reservation exists.
* `RESERVATION_NAME`: the name of the reservation.

The output is similar to the following:

```
{
  "items": [
    {
      "assuredCount": 1,
      "count": 1,
      "creationTimestamp": "2024-10-17T12:49:56.971-07:00",
      "id": "8544903383436022926",
      "inUseCount": 1,
      "kind": "compute#reservationBlock",
      "name": "example-res1-block-1",
      "physicalTopology": {
        "block": "c18707ac3d2493381c9f01fa775c4a68",
        "cluster": "europe-west1-cluster-jfhb"
      },
      "reservationMaintenance": {
        "schedulingType": "GROUPED"
      },
      "selfLink": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-res1/reservationBlocks/example-res1-block-1",
      "selfLinkWithId": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-res1/reservationBlocks/8544903383436022926",
      "status": "READY",
      "zone": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b"
    },
    {
      "assuredCount": 2,
      "count": 2,
      "creationTimestamp": "2024-10-17T12:49:54.354-07:00",
      "id": "5787689015406548144",
      "inUseCount": 2,
      "kind": "compute#reservationBlock",
      "name": "example-res1-block-2",
      "physicalTopology": {
        "block": "a9b7f2e4c6d1853902b4f5a7d8e31c60",
        "cluster": "europe-west1-cluster-jfhb"
      },
      "reservationMaintenance": {
        "schedulingType": "GROUPED"
      },
      "selfLink": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-res1/reservationBlocks/example-res1-block-2",
      "selfLinkWithId": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b/reservations/example-res1/reservationBlocks/5787689015406548144",
      "status": "READY",
      "zone": "https://www.googleapis.com/compute/beta/projects/example-project/zones/europe-west1-b"
    },
    ...
  ]
}

```

In the `physicalTopology` field of a block, you can see the physical
location of the block as follows:

* **`cluster`**: the global name of the cluster.
* **`block`**: the organization-specific ID of the block in which VMs will
  be located.

What's next
-----------

* VM and cluster creation overview