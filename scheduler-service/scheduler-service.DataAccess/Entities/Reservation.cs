namespace scheduler_service.DataAccess.Entities;

public enum ReservationStatus
{
    PendingConfirmation = 1,
    Confirmed = 2,
    Canceled = 3,
}

public class Reservation
{
    public long Id { get; set; }

    public long TimeslotId { get; set; }

    public virtual Timeslot Timeslot { get; set; } = null!;

    public ReservationStatus ReservationStatus { get; set; }
}